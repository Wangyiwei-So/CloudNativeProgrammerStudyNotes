# 实现类似go语言channel的优雅关闭

简介: 利用tokio的broadcast，实现服务的优雅关闭。即接收到ctrl+c命令后，主线程发送广播事件，通知其他以协程方式启动的服务优雅关闭，并在第二次ctrl+c后强制关闭。

首先介绍三个包

async_trait: 可以实现让trait中的方法带async

futures_util: 是一个为 Rust 异步编程提供基础的官方维护的库，以在 Rust 中进行零成本异步编程。详见[https://www.notion.so/futures-util-fc1dbaea5ff44cbabe7107b8b180fe69?pvs=4](https://www.notion.so/futures-util-fc1dbaea5ff44cbabe7107b8b180fe69?pvs=21)

anyhow: 用于统一处理错误

在[dependencies]中添加依赖

```toml
tokio = { version = "1.34.0", features = ["full"] }
async-trait = "0.1.74"
futures-util = "0.3.29"
anyhow = "1.0.75"
```

```rust
use std::pin::Pin;
use async_trait::async_trait;
use tokio::{sync::broadcast,time::{sleep, Duration}, signal::ctrl_c};
use futures_util::{FutureExt, Future, future::join_all};

/// 定义一个表示Service的trait,方便对服务进行统一的注册和管理
#[async_trait]
pub trait Service{
    async fn run(&mut self, ctx: broadcast::Receiver<()>)->anyhow::Result<()>;
}

/// 一个service,会每1秒打印Hello
struct PrintHelloService{}

#[async_trait]
impl Service for PrintHelloService{
    async fn run(&mut self, mut ctx: broadcast::Receiver<()>)->anyhow::Result<()>{
        println!("PrintHelloService starting");
        loop {
            tokio::select! { //类似go的写法
              _ = ctx.recv()=>{
                graceful_shutdown(|| Box::pin(async {   //传入闭包，用于执行优雅关闭的逻辑
                    println!("PrintHelloService 优雅关闭开始");
                    sleep(Duration::from_secs(6)).await;
                    println!("PrintHelloService 优雅关闭完成");
                })).await;
                return Ok(())
              }
              _ = sleep(Duration::from_secs(1)) => {    //注意sleep不能用标准库的，要用tokio的
                println!("Hello")
              }
            }
        }
    }
}

struct PrintWorldService{}

#[async_trait]
impl Service for PrintWorldService{
    async fn run(&mut self, mut ctx: broadcast::Receiver<()>)->anyhow::Result<()>{
        println!("PrintWorldService starting");
        loop {
            tokio::select! {
              _ = ctx.recv()=>{
                graceful_shutdown(|| Box::pin(async {
                    println!("PrintWorldService 优雅关闭开始");
                    sleep(Duration::from_secs(3)).await;
                    println!("PrintWorldService 优雅关闭完成");
                })).await;
                
                return Ok(())
              }
              _ = sleep(Duration::from_secs(1)) => {
                println!("World")
              }
            }
        }
    }
}

async fn graceful_shutdown<F>(callback: F)
where F: FnOnce()->Pin<Box<dyn Future<Output = ()> + Send>>{  //这个返回类型是futures_core::future::BoxFuture的类型，也就是带async的函数.boxed()的返回类型
    
    let (sender,_) = broadcast::channel(1);
    let shutdown_obj = start_signal(sender).boxed(); //注意.boxed()这个方法必须要use FutureExt才有
    tokio::select! {
        _ = callback() => {
            return
        }
        _ = shutdown_obj => {
            println!("[warn] 没有优雅关闭完成就强制停止了!")
        }
    }
}

async fn start_signal(sender: broadcast::Sender<()>)->anyhow::Result<()>{
    ctrl_c().await?;    //阻塞并接收信号
    println!("接收到ctrlc信号");
    sender.send(())?;
    Ok(())
}

/// 随便搞一个服务管理器
struct ServiceManager{
    services: Vec<Box<dyn Service>>,
}

#[tokio::main]
async fn main(){
    let mut sm = ServiceManager{services:vec![]};
    let service1 = PrintHelloService{};
    let service2 = PrintWorldService{};
    sm.services.push(Box::new(service1));
    sm.services.push(Box::new(service2));
    let (sender,ctx) = broadcast::channel(1);
    let ctx_obj = start_signal(sender).boxed();
    let mut services: Vec<_> = sm.services.iter_mut().map(|service|{
        service.run(ctx.resubscribe()).boxed()
    }).collect();
    services.push(ctx_obj);
    join_all(services).await;
}
```

测试代码，不用管

```rust
use anyhow::Result;
use std::time::Duration;
use futures_util::FutureExt;
use tokio_stream::{self as stream, StreamExt};
use tokio::time::sleep;
use tokio::signal::ctrl_c;
use tokio::sync::broadcast;

#[tokio::main(flavor = "multi_thread")]
async fn main() {
    start().await;
}

async fn start_signal(tx: broadcast::Sender<()>)->Result<()>{
    ctrl_c().await.unwrap();
    println!("优雅关闭");
    tx.send(());
    Ok(())
}

async fn start(){
    let (ctx,_) = broadcast::channel::<()>(1);
    let tick1_obj = tick1(ctx.clone()).fuse().boxed();
    let tick2_obj = tick2(ctx.clone()).fuse().boxed();
    let more_async_work_obj = more_async_work().fuse().boxed();
    let start_signal_obj = start_signal(ctx.clone()).fuse().boxed();
    // let services = Box::pin(async {
    //     tokio::select! {
    //         _ = tick1_obj => {},
    //         _ = tick2_obj => {},
    //     }
    //     Ok::<(), anyhow::Error>(())
    // });
    
    let _ = tokio::try_join!(tick1_obj,tick2_obj, more_async_work_obj, start_signal_obj);
}

async fn more_async_work() -> Result<()> {
    // more here
    println!("more_async_work");
    Ok(())
}

async fn tick1(tx: broadcast::Sender<()>)->Result<()>{
    let mut rx = tx.subscribe();
    for i in 0..5{
        tokio::select! {
            _ = rx.recv()=>{
                println!("tick1 优雅关闭");
                let (ctx,_) = broadcast::channel::<()>(1);
                let start_signal_obj = start_signal(ctx.clone()).fuse().boxed();
                tokio::select!{
                    _ = start_signal_obj => {
                        println!("中断优雅关闭");
                        return Ok(());
                    }
                    _ = sleep(Duration::from_secs(5)) => {
                        println!("tick1 优雅关闭结束");
                    }
                }
                break;
            }
            _ = sleep(Duration::from_secs(1)) => {
                println!("tick1 {}",i);
            }
        }
    }
    Ok(())
}

async fn tick2(tx: broadcast::Sender<()>)->Result<()>{
    let mut rx = tx.subscribe();
    for i in 0..5{
        tokio::select! {
            _ = rx.recv()=>{
                println!("tick2 优雅关闭");
                break;
            }
            _ = sleep(Duration::from_secs(1)) => {
                println!("tick2 {}",i);
            }
        }
    }
    Ok(())
}
```
