import React from 'react';

interface Props {
    uploader: string;
    campaign: string;
    onChangeUploader: (value: string) => void;
    onChangeCampaign: (value: string) => void;
    onSaveSettings: () => void;
}

const ConfigPanel = ({ uploader, campaign, onChangeUploader, onChangeCampaign, onSaveSettings }: Props) => (
    <section className="card">
        <h3 className="card-title">Configuration</h3>
        <div className="input-group">
            <label className="label">Uploader Name</label>
            <input
                className="input-field"
                value={uploader}
                onChange={(e) => onChangeUploader(e.target.value)}
                onBlur={onSaveSettings}
            />
        </div>
        <div className="input-group" style={{ marginTop: '12px' }}>
            <label className="label">Campaign ID</label>
            <input
                className="input-field"
                value={campaign}
                onChange={(e) => onChangeCampaign(e.target.value)}
                onBlur={onSaveSettings}
            />
        </div>
    </section>
);

export default ConfigPanel;
