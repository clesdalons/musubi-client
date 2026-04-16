import { useState, useEffect } from 'react';
import { appApi } from '../services/appApi';

export function useWatcher() {
    const [path, setPath] = useState<string>('');
    const [status, setStatus] = useState<string>('Initializing...');
    const [info, setInfo] = useState({ name: 'Loading...', version: '' });

    useEffect(() => {
        appApi.getAppInfo().then((appInfo) => setInfo(appInfo));

        appApi.getSettings().then((result) => {
            setPath(result.save_path);
            setStatus(result.save_path ? 'Watching' : 'Setup Required');
        });

        appApi.onEvent('watcher:detected', () => {
            setStatus('New save detected!');
            setTimeout(() => setStatus('Watching'), 3000);
        });

        return () => {
            appApi.offEvent('watcher:detected');
        };
    }, []);

    return { path, setPath, status, setStatus, info };
}
