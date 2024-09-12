use futures_util::StreamExt;
use reqwest::Client;
use serde::Deserialize;
use std::fs;
use std::io;
use std::io::Write;
use std::process::Command;
use sysinfo::System;
use tokio::io::AsyncWriteExt;

// Why updater using Rust? mean while the main application is written in Go?
//
// Im in halfway creating EmuWatcher using Golang, but i got memeory isssue
// but changing to Rust, it will took more day to finish it, so i decided to
// keep the main application in Go and create the updater in Rust
// Also both language is good and im learning them

#[derive(Deserialize)]
struct ReleaseAsset {
    browser_download_url: String,
    name: String,
}

#[derive(Deserialize)]
struct Release {
    tag_name: String,
    assets: Vec<ReleaseAsset>,
}

#[tokio::main]
async fn main() {
    close_emuwatcher_processes();
    match get_data().await {
        Ok((version, download_url)) => match get_current_version() {
            Ok(current_version) => {
                if current_version.trim() == version.trim() {
                    println!("You are already using the latest version");
                    println!("Press Enter to exit and open program...");
                    let mut input = String::new();
                    let _ = io::stdin().read_line(&mut input);
                } else {
                    println!("New version available: {}", version);
                    println!("Downloading from: {}", download_url);
                    if let Err(e) = download(download_url.trim()).await {
                        eprintln!("Download failed: {}", e);
                    }
                }
            }
            Err(e) => {
                eprintln!("Error getting current version: {}", e);
                println!("Downloading latest version...");
                if let Err(e) = download(download_url.trim()).await {
                    eprintln!("Download failed: {}", e);
                }
            }
        },
        Err(e) => eprintln!("Error: {}", e),
    }
}

fn close_emuwatcher_processes() {
    let mut system = System::new_all();
    system.refresh_all();

    for (_pid, process) in system.processes() {
        if process.name().to_str().unwrap().to_lowercase() == "emuwatcher.exe" {
            println!("Closing EmuWatcher.exe");
            process.kill();
        }
    }
}

async fn download(url: &str) -> Result<(), Box<dyn std::error::Error>> {
    let client = Client::new();
    let response = client.get(url).send().await?;

    let total_size = response.content_length().unwrap_or(0);

    let mut file = tokio::fs::File::create("EmuWatcher.7z").await?;
    let mut stream = response.bytes_stream();

    let mut downloaded: u64 = 0;

    while let Some(chunk) = stream.next().await {
        let chunk = chunk?;
        downloaded += chunk.len() as u64;
        file.write_all(&chunk).await?;
        print!(
            "\rDownloading... {:.2}%",
            (downloaded as f64 / total_size as f64) * 100.0
        );
        io::stdout().flush()?;
    }
    print!("\n{:<width$}\r", "Extracting...", width = 30);
    io::stdout().flush()?;
    sevenz_rust::decompress_file("EmuWatcher.7z", "./").expect("complete");
    std::fs::remove_file("EmuWatcher.7z").expect("complete");
    print!("{:<width$}\r", "Completed!\n", width = 30);
    println!("Press Enter to exit and open program...");
    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    let _ = Command::new("cmd")
        .arg("/c")
        .arg("start")
        .arg("EmuWatcher.exe")
        .spawn();
    Ok(())
}

fn get_current_version() -> Result<String, Box<dyn std::error::Error>> {
    let emu_file = "EmuWatcher.exe";

    if !fs::metadata(emu_file).is_ok() {
        return Err("EmuWatcher not found".into());
    }

    let output = Command::new(emu_file).arg("version").output()?;
    let version = String::from_utf8(output.stdout)?;

    if version.trim().is_empty() {
        Err("Version is empty".into())
    } else {
        Ok(version.trim().to_string())
    }
}

async fn get_data() -> Result<(String, String), Box<dyn std::error::Error>> {
    let repo = "Irilith/EmuWatcher";
    let url = format!("https://api.github.com/repos/{}/releases/latest", repo);

    let client = Client::new();
    let response = client
        .get(&url)
        .header("User-Agent", "request")
        .send()
        .await?;

    if response.status().is_success() {
        let release: Release = response.json().await?;
        let version = release.tag_name;

        if let Some(asset) = release.assets.iter().find(|a| a.name == "EmuWatcher.7z") {
            let download_url = asset.browser_download_url.clone();
            Ok((version, download_url))
        } else {
            Err("No matching asset found".into())
        }
    } else {
        Err(format!("Failed to fetch the version. Status: {}", response.status()).into())
    }
}
