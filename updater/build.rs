fn main() {
    if std::env::var("CARGO_CFG_TARGET_OS").unwrap() == "windows" {
        let mut res = winresource::WindowsResource::new();
        res.set_manifest_file("../assets/Updater.exe.xml");
        res.set_icon("../assets/Updater_Logo.ico");
        res.compile().unwrap();
    }
}
