import { useState, useEffect } from 'react';
import { useWatcher } from './hooks/useWatcher';
import { SelectFolder, GetSettings, SaveSettings } from "../wailsjs/go/main/App";
import './App.css';

const App = () => {
    const { path, setPath, lastSave, status, setStatus, info } = useWatcher();
    const [uploader, setUploader] = useState("");
    const [campaign, setCampaign] = useState("");

    useEffect(() => {
        GetSettings().then(cfg => {
            setUploader(cfg.uploader || "");
            setCampaign(cfg.campaign || "");
        });
    }, []);

    const handleSyncSettings = () => {
        SaveSettings(path, uploader, campaign);
    };

    const handleBrowse = async () => {
        const selected = await SelectFolder();
        if (selected) {
            setPath(selected);
            setStatus("Watching");
        }
    };

    return (
        <div className="container">
            <header className="header">
                <h1 className="title">{info.name}</h1>
                <div className="badge">{status}</div>
            </header>

            <main className="main">
                <section className="card">
                    <h3 className="card-title">Configuration</h3>
                    <div className="input-group">
                        <label className="label">Uploader Name</label>
                        <input 
                            className="input-field"
                            value={uploader}
                            onChange={(e) => setUploader(e.target.value)}
                            onBlur={handleSyncSettings}
                            placeholder="e.g. Incurso"
                        />
                    </div>
                    <div className="input-group" style={{ marginTop: '12px' }}>
                        <label className="label">Campaign ID</label>
                        <input 
                            className="input-field"
                            value={campaign}
                            onChange={(e) => setCampaign(e.target.value)}
                            onBlur={handleSyncSettings}
                            placeholder="e.g. testcampaign"
                        />
                    </div>
                </section>

                <section className="card">
                    <h3 className="card-title">Save Directory</h3>
                    <div className="path-box">
                        <code className="code-block">{path || "No folder selected"}</code>
                        <button onClick={handleBrowse} className="btn-mini">
                            {path ? "Change" : "Browse"}
                        </button>
                    </div>
                </section>

                <section className="card">
                    <h3 className="card-title">Status</h3>
                    <div className="save-display">
                        <span className="label">Latest synced:</span>
                        <span className="file-name">{lastSave}</span>
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