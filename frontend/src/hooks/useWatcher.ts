import { useState, useEffect } from 'react';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { GetAppInfo, GetInitialPath } from '../../wailsjs/go/application/App';

export function useWatcher() {
    const [path, setPath] = useState<string>('');
    const [status, setStatus] = useState<string>('Initializing...');
    const [info, setInfo] = useState({ name: 'Loading...', version: '' });

    useEffect(() => {
        GetAppInfo().then((appInfo) => setInfo(appInfo));

        GetInitialPath().then((result) => {
            setPath(result);
            setStatus(result ? 'Watching' : 'Setup Required');
        });

        EventsOn('watcher:detected', () => {
            setStatus('New save detected!');
            setTimeout(() => setStatus('Watching'), 3000);
        });

        return () => {
            EventsOff('watcher:detected');
        };
    }, []);

    return { path, setPath, status, setStatus, info };
}
