import { useState, useEffect } from 'react';
import { GetSavePath } from "../wailsjs/go/main/App"; // L'import magique
import './App.css';

function App() {
    const [path, setPath] = useState("Recherche...");

    useEffect(() => {
        // On appelle la fonction Go au démarrage
        GetSavePath().then((result) => {
            setPath(result);
        });
    }, []);

    return (
        <div id="App">
            <header style={{ padding: '20px', textAlign: 'center' }}>
                <h1 style={{ color: '#61dafb' }}>Musubi</h1>
            </header>
            <main style={{ padding: '0 20px' }}>
                <div className="card">
                    <p>Dossier détecté :</p>
                    <code style={{ fontSize: '0.8em', wordBreak: 'break-all' }}>
                        {path}
                    </code>
                </div>
                <button style={{ marginTop: '20px' }}>
                    Synchroniser
                </button>
            </main>
        </div>
    );
}

export default App;