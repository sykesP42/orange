import { Camera, CameraResultType, CameraSource } from '@capacitor/camera'
import { Geolocation, Position } from '@capacitor/geolocation'
import { Network } from '@capacitor/network'
import { LocalNotifications } from '@capacitor/local-notifications'
import { Share } from '@capacitor/share'
import { Preferences } from '@capacitor/preferences'

export class MobileService {
  private static instance: MobileService
  private isNative: boolean = false

  static getInstance(): MobileService {
    if (!MobileService.instance) {
      MobileService.instance = new MobileService()
    }
    return MobileService.instance
  }

  constructor() {
    this.isNative = this.checkIfNative()
    if (this.isNative) {
      this.initServices()
    }
  }

  private checkIfNative(): boolean {
    return typeof (window as any).Capacitor !== 'undefined'
  }

  private async initServices() {
    await LocalNotifications.requestPermissions()
    
    const network = await Network.getStatus()
    console.log('Network status:', network.connected)
    
    Network.addListener('networkStatusChange', (status) => {
      window.dispatchEvent(new CustomEvent('network-change', { detail: status }))
    })
  }

  async takePhoto(options?: { quality?: number; allowEditing?: boolean }): Promise<string | null> {
    try {
      if (!this.isNative) {
        return null
      }

      const photo = await Camera.getPhoto({
        quality: options?.quality || 90,
        allowEditing: options?.allowEditing || false,
        resultType: CameraResultType.Base64,
        source: CameraSource.Prompt
      })

      return photo.base64String ? `data:image/jpeg;base64,${photo.base64String}` : null
    } catch (error) {
      console.error('Camera error:', error)
      return null
    }
  }

  async pickFromGallery(): Promise<string | null> {
    try {
      if (!this.isNative) {
        return null
      }

      const photo = await Camera.getPhoto({
        quality: 90,
        resultType: CameraResultType.Base64,
        source: CameraSource.Photos
      })

      return photo.base64String ? `data:image/jpeg;base64,${photo.base64String}` : null
    } catch (error) {
      console.error('Gallery picker error:', error)
      return null
    }
  }

  async getCurrentPosition(): Promise<Position | null> {
    try {
      if (!this.isNative) {
        return null
      }

      const coordinates = await Geolocation.getCurrentPosition({
        enableHighAccuracy: true,
        timeout: 10000
      })
      
      return coordinates
    } catch (error) {
      console.error('Geolocation error:', error)
      return null
    }
  }

  async isOnline(): Promise<boolean> {
    if (!this.isNative) {
      return navigator.onLine
    }

    const status = await Network.getStatus()
    return status.connected
  }

  async sendNotification(title: string, body: string, id?: number): Promise<void> {
    if (!this.isNative) {
      if ('Notification' in window && Notification.permission === 'granted') {
        new Notification(title, { body })
      }
      return
    }

    await LocalNotifications.schedule({
      notifications: [{
        title,
        body,
        id: id || Date.now(),
        schedule: { at: new Date(Date.now() + 1000) },
        sound: 'default',
        smallIcon: 'ic_stat_icon',
        attachments: []
      }]
    })
  }

  async shareContent(options: { title: string; text?: string; url?: string; dialogTitle?: string }): Promise<void> {
    if (!this.isNative) {
      if (navigator.share) {
        await navigator.share(options)
      }
      return
    }

    await Share.share(options)
  }

  async savePreference(key: string, value: string): Promise<void> {
    await Preferences.set({ key, value })
  }

  async getPreference(key: string): Promise<string | null> {
    const { value } = await Preferences.get({ key })
    return value
  }

  async removePreference(key: string): Promise<void> {
    await Preferences.remove({ key })
  }

  getIsNative(): boolean {
    return this.isNative
  }
}

export const mobileService = MobileService.getInstance()
