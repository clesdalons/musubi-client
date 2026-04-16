import React from 'react';

interface Props {
    path: string;
    onBrowse: () => void;
    onOpen: () => void;
}

const DirectoryPanel = ({ path, onBrowse, onOpen }: Props) => (
    <section className="card">
        <h3 className="card-title">Save Directory</h3>
        <div className="path-box">
            <code className="code-block">{path || 'No folder selected'}</code>
            <div className="path-actions">
                <button onClick={onOpen} className="btn-mini" disabled={!path}>
                    Open
                </button>
                <button onClick={onBrowse} className="btn-mini">
                    Change
                </button>
            </div>
        </div>
    </section>
);

export default DirectoryPanel;
