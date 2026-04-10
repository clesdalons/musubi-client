import { useState, useEffect } from 'react';
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { GetInitialPath, StartWatcher } from "../../wailsjs/go/main/App";
import { GetAppInfo } from "../../wailsjs/go/main/App";

export function useWatcher() {
    const [path, setPath] = useState<string>("");
    const [lastSave, setLastSave] = useState<string>("Waiting for changes...");
    const [status, setStatus] = useState<string>("Initializing...");
    const [info, setInfo] = useState({ name: "Loading...", version: "" });

    useEffect(() => {
        GetAppInfo().then(appInfo => setInfo(appInfo));
        // Initialize path and watcher
        GetInitialPath().then((result) => {
            setPath(result);
            if (result) {
                StartWatcher();
                setStatus("Watching");
            } else {
                setStatus("Setup Required");
            }
        });

        // Listen for Go events
        const unsubscribe = EventsOn("new-save-event", (fileName: string) => {
            setLastSave(fileName);
            setStatus("New save detected!");
            setTimeout(() => setStatus("Watching"), 3000);
        });

        return () => unsubscribe();
    }, []);

    return { path, setPath, lastSave, status, setStatus, info };
}