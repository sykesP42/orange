import { invoke } from '@tauri-apps/api/core'
import { save, open } from '@tauri-apps/plugin-dialog'
import { sendNotification } from '@tauri-apps/plugin-notification'

export class DesktopService {
  private static instance: DesktopService
  private isTauri: boolean = false

  static getInstance(): DesktopService {
    if (!DesktopService.instance) {
      DesktopService.instance = new DesktopService()
    }
    return DesktopService.instance
  }

  constructor() {
    this.isTauri = this.checkIfTauri()
  }

  private checkIfTauri(): boolean {
    return '__TAURI__' in window
  }

  getIsTauri(): boolean {
    return this.isTauri
  }

  async exportToExcel(data: any[], filename?: string): Promise<string | null> {
    if (!this.isTauri) {
      console.warn('Export to Excel is only available in Tauri desktop app')
      return null
    }

    try {
      const filePath = await save({
        defaultPath: filename || 'bloggers_export.xlsx',
        filters: [
          { name: 'Excel Files', extensions: ['xlsx', 'xls'] }
        ]
      })

      if (!filePath) return null

      await invoke('export_to_excel', { data, path: filePath })
      return filePath
    } catch (error) {
      console.error('Export to Excel error:', error)
      return null
    }
  }

  async exportToCSV(data: any[], filename?: string): Promise<string | null> {
    if (!this.isTauri) {
      console.warn('Export to CSV is only available in Tauri desktop app')
      return null
    }

    try {
      const filePath = await save({
        defaultPath: filename || 'bloggers_export.csv',
        filters: [
          { name: 'CSV Files', extensions: ['csv'] }
        ]
      })

      if (!filePath) return null

      await invoke('export_to_csv', { data, path: filePath })
      return filePath
    } catch (error) {
      console.error('Export to CSV error:', error)
      return null
    }
  }

  async importFromFile(): Promise<string | null> {
    if (!this.isTauri) {
      console.warn('Import file is only available in Tauri desktop app')
      return null
    }

    try {
      const selectedPath = await open({
        multiple: false,
        filters: [
          { name: 'Excel/CSV', extensions: ['xlsx', 'xls', 'csv'] },
          { name: 'JSON', extensions: ['json'] },
          { name: 'All Files', extensions: ['*'] }
        ]
      })

      if (!selectedPath) return null

      const content = await invoke('read_local_file', { path: selectedPath })
      return content as string
    } catch (error) {
      console.error('Import file error:', error)
      return null
    }
  }

  async getSystemInfo(): Promise<{ os: string; arch: string; family: string } | null> {
    if (!this.isTauri) return null

    try {
      const info = await invoke('get_system_info')
      return info as { os: string; arch: string; family: string }
    } catch (error) {
      console.error('Get system info error:', error)
      return null
    }
  }

  async showNotification(title: string, body: string): Promise<void> {
    if (!this.isTauri) {
      if ('Notification' in window && Notification.permission === 'granted') {
        new Notification(title, { body })
      }
      return
    }

    try {
      sendNotification({ title, body })
    } catch (error) {
      console.error('Show notification error:', error)
    }
  }

  minimizeWindow(): void {
    if (this.isTauri) {
      import('@tauri-apps/api/window').then(({ getCurrentWindow }) => {
        getCurrentWindow().minimize()
      })
    }
  }

  maximizeWindow(): void {
    if (this.isTauri) {
      import('@tauri-apps/api/window').then(({ getCurrentWindow }) => {
        getCurrentWindow().toggleMaximize()
      })
    }
  }

  closeWindow(): void {
    if (this.isTauri) {
      import('@tauri-apps/api/window').then(({ getCurrentWindow }) => {
        getCurrentWindow().close()
      })
    }
  }
}

export const desktopService = DesktopService.getInstance()
