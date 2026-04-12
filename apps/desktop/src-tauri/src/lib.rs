use tauri::Manager;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_notification::init())
        .plugin(tauri_plugin_shell::init())
        .setup(|app| {
            #[cfg(debug_assertions)]
            {
                let window = app.get_webview_window("main").unwrap();
                window.open_devtools();
            }
            
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            export_to_excel,
            export_to_csv,
            read_local_file,
            get_system_info
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[tauri::command]
async fn export_to_excel(data: Vec<BloggerRecord>, path: String) -> Result<(), String> {
    use calamine::{Writer, Xlsx};
    
    let mut workbook = Xlsx::new();
    let mut sheet = workbook.worksheet("Bloggers").map_err(|e| e.to_string())?;
    
    // Header row
    sheet.write_row(&["ID", "Nickname", "Category", "Platform", "Status", "Contact", "Created At"]).map_err(|e| e.to_string())?;
    
    // Data rows
    for record in data {
        sheet.write_row(&[
            record.id,
            record.nickname,
            record.category,
            record.platform,
            record.status,
            record.contact,
            record.created_at
        ]).map_err(|e| e.to_string())?;
    }
    
    workbook.save(&path).map_err(|e| e.to_string())?;
    Ok(())
}

#[tauri::command]
fn export_to_csv(data: Vec<BloggerRecord>, path: String) -> Result<(), String> {
    use std::fs::File;
    use std::io::Write;
    
    let mut file = File::create(&path).map_err(|e| e.to_string())?;
    
    // Header
    writeln!(file, "ID,Nickname,Category,Platform,Status,Contact,Created At").map_err(|e| e.to_string())?;
    
    // Data
    for record in data {
        writeln!(file, "{},{},{},{},{},{},{}", 
            record.id, 
            record.nickname, 
            record.category, 
            record.platform, 
            record.status, 
            record.contact, 
            record.created_at
        ).map_err(|e| e.to_string())?;
    }
    
    Ok(())
}

#[tauri::command]
async fn read_local_file(path: String) -> Result<String, String> {
    use std::fs;
    
    fs::read_to_string(&path).map_err(|e| e.to_string())
}

#[tauri::command]
fn get_system_info() -> SystemInfo {
    SystemInfo {
        os: std::env::consts::OS.to_string(),
        arch: std::env::consts::ARCH.to_string(),
        family: std::env::consts::FAMILY.to_string(),
    }
}

#[derive(serde::Serialize)]
struct BloggerRecord {
    id: i32,
    nickname: String,
    category: String,
    platform: String,
    status: String,
    contact: String,
    created_at: String,
}

#[derive(serde::Serialize)]
struct SystemInfo {
    os: String,
    arch: String,
    family: String,
}
