import React from 'react';
import { CloudStatus, PullStatus } from '../types';
import dayjs from 'dayjs';

interface Props {
    localDate: string | null;
    cloudData: CloudStatus | null;
    lastCheck: dayjs.Dayjs | null;
    pullStatus: PullStatus;
    isRefreshing: boolean;
    onRefresh: () => void;
    onPull: () => void;
}

const SyncPanel = ({ localDate, cloudData, lastCheck, pullStatus, isRefreshing, onRefresh, onPull }: Props) => (
    <section className="card">
        <div className="card-header-flex">
            <h3 className="card-title">Cloud Synchronization</h3>
            <span className="last-check-text">
                {lastCheck ? `Updated ${lastCheck.fromNow()}` : 'Not checked'}
            </span>
        </div>

        <div className="sync-grid">
            <div className="sync-item">
                <span className="sync-label">Local Save</span>
                <span className="sync-value">{localDate ? dayjs(localDate).format('DD MMM HH:mm') : 'None'}</span>
            </div>
            <div className="sync-item">
                <span className="sync-label">Cloud Save</span>
                <span className="sync-value">
                    {cloudData?.timestamp ? dayjs(cloudData.timestamp as string).format('DD MMM HH:mm') : 'None'}
                </span>
                {cloudData?.uploader && <span className="sync-subvalue">by {cloudData.uploader}</span>}
            </div>
        </div>

        <div className="sync-actions">
            <button className="btn-mini" onClick={onRefresh} disabled={isRefreshing || pullStatus === 'loading'}>
                {isRefreshing ? 'Refreshing...' : 'Refresh Status'}
            </button>
            <button className={`btn-sync ${pullStatus}`} onClick={onPull} disabled={pullStatus !== 'idle' || isRefreshing}>
                {pullStatus === 'idle' && 'Cloud Pull'}
                {pullStatus === 'loading' && 'Pulling...'}
                {pullStatus === 'success' && '✓ Success'}
                {pullStatus === 'error' && '✕ Failed'}
            </button>
        </div>
    </section>
);

export default SyncPanel;
