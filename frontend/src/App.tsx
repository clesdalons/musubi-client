import { useState, useEffect } from 'react';
import { useWatcher } from './hooks/useWatcher';
import { SelectFolder, GetSettings, SaveSettings, DownloadLatestSave } from "../wailsjs/go/main/App";
import './App.css';

const App = () => {
    const { path, setPath, lastSave, status, setStatus, info } = useWatcher();
    const [uploader, setUploader] = useState("");
    const [campaign, setCampaign] = useState("");
    const [isDownloading, setIsDownloading] = useState(false); // État pour le bouton Pull

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

    const [pullStatus, setPullStatus] = useState<"idle" | "loading" | "success" | "error">("idle");

    const handleCloudPull = async () => {
        setPullStatus("loading");
        setStatus("Downloading...");
        
        try {
            const result = await DownloadLatestSave();
            
            if (result === "Success") {
                setPullStatus("success");
                setStatus("Watching");

                // Le bouton revient à la normale après 3 secondes
                setTimeout(() => {
                    setPullStatus("idle");
                }, 3000);
            } else {
                setPullStatus("error");
                console.error(result);
                setTimeout(() => setPullStatus("idle"), 5000);
            }
        } catch (err) {
            setPullStatus("error");
            setTimeout(() => setPullStatus("idle"), 5000);
        }
    };

    return (
        <div className="container">
            <header className="header">
                <h1 className="title">{info.name}</h1>
                <div className="badge">{status}</div>
            </header>

            <main className="main">
                {/* Section Configuration */}
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

                {/* Section Save Directory & Cloud Pull */}
                <section className="card">
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
                        <h3 className="card-title" style={{ margin: 0 }}>Save Directory</h3>
                        {/* Nouveau bouton Cloud Pull */}
                        <button 
                            onClick={handleCloudPull} 
                            disabled={pullStatus !== "idle" || !path}
                            className={`btn-sync ${pullStatus}`}
                        >
                            {pullStatus === "idle" && "Cloud Pull"}
                            {pullStatus === "loading" && "Pulling..."}
                            {pullStatus === "success" && "✓ Success"}
                            {pullStatus === "error" && "✕ Failed"}
                        </button>
                    </div>
                    
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
                        <span className="file-name">{lastSave || "None"}</span>
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