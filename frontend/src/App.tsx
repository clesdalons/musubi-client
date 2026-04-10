import { useWatcher } from './hooks/useWatcher';
import { SelectFolder } from "../wailsjs/go/main/App";
import './App.css';

const App = () => {
    const { path, setPath, lastSave, status, setStatus, info } = useWatcher();

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
                    <h3 className="card-title">Save Directory</h3>
                    {path ? (
                        <div className="path-box">
                            <code className="code-block">{path}</code>
                            <button onClick={handleBrowse} className="btn-mini">Change Path</button>
                        </div>
                    ) : (
                        <div style={{ textAlign: 'center' }}>
                            <button onClick={handleBrowse} className="btn-primary">
                                Browse Folder
                            </button>
                        </div>
                    )}
                </section>

                <section className="card">
                    <h3 className="card-title">Live Monitoring</h3>
                    <div className="save-display">
                        <span className="label">Latest sync-ready file:</span>
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