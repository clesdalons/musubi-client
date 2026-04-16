import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import {
    GetAppInfo as getAppInfoNative,
    GetSettings as getSettingsNative,
    SaveSettings as saveSettingsNative,
    SelectFolder as selectFolderNative,
    OpenFolder as openFolderNative,
    GetLocalSaveStatus as getLocalSaveStatusNative,
    GetCloudSaveStatus as getCloudSaveStatusNative,
    DownloadLatestSave as downloadLatestSaveNative,
} from '../../wailsjs/go/application/App';
import { AppInfo, Config, CloudStatus } from '../types';

export type WailsEventName = 'watcher:detected' | 'upload:success' | 'upload:error';

// Application API boundary.
// This module keeps the rest of the UI platform-agnostic
// and centralizes Wails integration in one place.
export const appApi = {
    getAppInfo: async (): Promise<AppInfo> => getAppInfoNative(),
    getSettings: async (): Promise<Config> => getSettingsNative(),
    saveSettings: async (path: string, uploader: string, campaign: string): Promise<void> => {
        await saveSettingsNative(path, uploader, campaign);
    },
    selectFolder: async (): Promise<string> => selectFolderNative(),
    openFolder: async (path: string): Promise<void> => {
        await openFolderNative(path);
    },
    getLocalSaveStatus: async (): Promise<string> => getLocalSaveStatusNative(),
    getCloudSaveStatus: async (): Promise<CloudStatus> => getCloudSaveStatusNative(),
    downloadLatestSave: async (): Promise<string> => downloadLatestSaveNative(),
    onEvent: (name: WailsEventName, callback: () => void): (() => void) => EventsOn(name, callback),
    offEvent: (name: WailsEventName): void => EventsOff(name),
};
