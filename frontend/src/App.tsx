import { useState, useEffect } from 'react';
import { useWatcher } from './hooks/useWatcher';
import { useSync } from './hooks/useSync';
import { SelectFolder, GetSettings, SaveSettings, OpenFolder } from "../wailsjs/go/main/App";
import { EventsOn, EventsOff } from "../wailsjs/runtime/runtime";
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import './App.css';

dayjs.extend(relativeTime);

const App = () => {
    const { path, setPath, status, setStatus, info } = useWatcher();
    const [uploader, setUploader] = useState("");
    const [campaign, setCampaign] = useState("");
    
    // État local pour gérer la couleur/classe du badge (idle, detected, success, error)
    const [badgeClass, setBadgeClass] = useState("");

    const { localDate, cloudData, lastCheck, pullStatus, isRefreshing, checkStatus, performPull } = useSync(campaign, setStatus);

    useEffect(() => {
        const initApp = async () => {
            const cfg = await GetSettings();
            setUploader(cfg.uploader || "");
            setCampaign(cfg.campaign || "");
            setPath(cfg.save_path || "");

            if (cfg.campaign) {
                checkStatus(cfg.campaign);
            }
        };

        initApp();

        // --- Écouteurs d'événements pour l'upload ---
        EventsOn("watcher:detected", () => {
            setStatus("New Save Detected : Uploading...");
            setBadgeClass("detected");
        });

        EventsOn("upload:success", () => {
            setStatus("Upload Successful");
            setBadgeClass("success"); 
            checkStatus();
            setTimeout(() => {
                setStatus("Watching");
                setBadgeClass("");
            }, 3000);
        });

        EventsOn("upload:error", () => {
            setStatus("Upload Failed");
            setBadgeClass("error");    // Devra être Rouge en CSS
            setTimeout(() => {
                setStatus("Watching");
                setBadgeClass("");
            }, 5000);
        });

        return () => {
            EventsOff("watcher:detected");
            EventsOff("upload:success");
            EventsOff("upload:error");
        };
    }, [checkStatus, setStatus]);

    const handleSyncSettings = () => SaveSettings(path, uploader, campaign);

    const handleBrowse = async () => {
        const selected = await SelectFolder();
        if (selected) {
            setPath(selected);
            setStatus("Watching");
        }
    };

    const handleOpenFolder = () => {
        if (path) OpenFolder(path);
    };

    return (
        <div className="container">
            <header className="header">
                <h1 className="title">{info.name}</h1>
                {/* On ajoute la badgeClass ici */}
                <div className={`badge ${badgeClass}`}>{status}</div>
            </header>

            <main className="main">
                <section className="card">
                    <h3 className="card-title">Configuration</h3>
                    <div className="input-group">
                        <label className="label">Uploader Name</label>
                        <input className="input-field" value={uploader} onChange={(e) => setUploader(e.target.value)} onBlur={handleSyncSettings} />
                    </div>
                    <div className="input-group" style={{ marginTop: '12px' }}>
                        <label className="label">Campaign ID</label>
                        <input className="input-field" value={campaign} onChange={(e) => setCampaign(e.target.value)} onBlur={handleSyncSettings} />
                    </div>
                </section>

                <section className="card">
                    <h3 className="card-title">Save Directory</h3>
                    <div className="path-box">
                        <code className="code-block">{path || "No folder selected"}</code>
                        <div className="path-actions">
                            <button onClick={handleOpenFolder} className="btn-mini" disabled={!path}>
                                Open
                            </button>
                            <button onClick={handleBrowse} className="btn-mini">
                                Change
                            </button>
                        </div>
                    </div>
                </section>
                
                <section className="card">
                    <div className="card-header-flex">
                        <h3 className="card-title">Cloud Synchronization</h3>
                        <span className="last-check-text">
                            {lastCheck ? `Updated ${lastCheck.fromNow()}` : "Not checked"}
                        </span>
                    </div>

                    <div className="sync-grid">
                        <div className="sync-item">
                            <span className="sync-label">Local Save</span>
                            <span className="sync-value">{localDate ? dayjs(localDate).format('DD MMM HH:mm') : "None"}</span>
                        </div>
                        <div className="sync-item">
                            <span className="sync-label">Cloud Save</span>
                            <span className="sync-value">
                                {cloudData?.timestamp ? dayjs(cloudData.timestamp).format('DD MMM HH:mm') : "None"}
                            </span>
                            {cloudData?.uploader && <span className="sync-subvalue">by {cloudData.uploader}</span>}
                        </div>
                    </div>

                    <div className="sync-actions">
                        <button 
                            className="btn-mini" 
                            onClick={() => checkStatus()}
                            disabled={isRefreshing || pullStatus === "loading"}
                        >
                            {isRefreshing ? "Refreshing..." : "Refresh Status"}
                        </button>
                        
                        <button 
                            className={`btn-sync ${pullStatus}`} 
                            onClick={performPull}
                            disabled={pullStatus !== "idle" || isRefreshing || !path}
                        >
                            {pullStatus === "idle" && "Cloud Pull"}
                            {pullStatus === "loading" && "Pulling..."}
                            {pullStatus === "success" && "✓ Success"}
                            {pullStatus === "error" && "✕ Failed"}
                        </button>
                    </div>
                </section>               
            </main>

            <footer className="footer">
                <p>Musubi Client • v{info.version}</p>
            </footer>
        </div>
    );
};

export default App;