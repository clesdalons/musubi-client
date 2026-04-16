export interface AppInfo {
    name: string;
    version: string;
}

export interface Config {
    save_path: string;
    uploader: string;
    campaign: string;
}

export type PullStatus = "idle" | "loading" | "success" | "error";

export interface CloudStatus {
    timestamp?: string;
    uploader?: string;
    [key: string]: unknown;
}
