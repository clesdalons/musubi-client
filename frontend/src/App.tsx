import { useState, useEffect } from 'react';
import { useWatcher } from './hooks/useWatcher';
import { useSync } from './hooks/useSync';
import { SelectFolder, GetSettings, SaveSettings, OpenFolder } from '../wailsjs/go/application/App';
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import ConfigPanel from './components/ConfigPanel';
import DirectoryPanel from './components/DirectoryPanel';
import SyncPanel from './components/SyncPanel';
import './App.css';

dayjs.extend(relativeTime);

const App = () => {
    const { path, setPath, status, setStatus, info } = useWatcher();
    const [uploader, setUploader] = useState('');
    const [campaign, setCampaign] = useState('');
    const [badgeClass, setBadgeClass] = useState('');

    const { localDate, cloudData, lastCheck, pullStatus, isRefreshing, checkStatus, performPull } = useSync(campaign, setStatus);

    useEffect(() => {
        const initApp = async () => {
            const cfg = await GetSettings();
            setUploader(cfg.uploader || '');
            setCampaign(cfg.campaign || '');
            setPath(cfg.save_path || '');

            if (cfg.campaign) {
                checkStatus(cfg.campaign);
            }
        };

        initApp();

        EventsOn('watcher:detected', () => {
            setStatus('New save detected : Uploading...');
            setBadgeClass('detected');
        });

        EventsOn('upload:success', () => {
            setStatus('Upload Successful');
            setBadgeClass('success');
            checkStatus();
            setTimeout(() => {
                setStatus('Watching');
                setBadgeClass('');
            }, 3000);
        });

        EventsOn('upload:error', () => {
            setStatus('Upload Failed');
            setBadgeClass('error');
            setTimeout(() => {
                setStatus('Watching');
                setBadgeClass('');
            }, 5000);
        });

        return () => {
            EventsOff('watcher:detected');
            EventsOff('upload:success');
            EventsOff('upload:error');
        };
    }, [checkStatus, setPath]);

    const handleSyncSettings = async () => {
        await SaveSettings(path, uploader, campaign);
    };

    const handleBrowse = async () => {
        const selected = await SelectFolder();
        if (selected) {
            setPath(selected);
            setStatus('Watching');
        }
    };

    const handleOpenFolder = () => {
        if (path) OpenFolder(path);
    };

    return (
        <div className="container">
            <header className="header">
                <h1 className="title">{info.name}</h1>
                <div className={`badge ${badgeClass}`}>{status}</div>
            </header>

            <main className="main">
                <ConfigPanel
                    uploader={uploader}
                    campaign={campaign}
                    onChangeUploader={setUploader}
                    onChangeCampaign={setCampaign}
                    onSaveSettings={handleSyncSettings}
                />

                <DirectoryPanel path={path} onBrowse={handleBrowse} onOpen={handleOpenFolder} />

                <SyncPanel
                    localDate={localDate}
                    cloudData={cloudData}
                    lastCheck={lastCheck}
                    pullStatus={pullStatus}
                    isRefreshing={isRefreshing}
                    onRefresh={() => checkStatus()}
                    onPull={performPull}
                />
            </main>

            <footer className="footer">
                <p>Musubi Client � v{info.version}</p>
            </footer>
        </div>
    );
};

export default App;
