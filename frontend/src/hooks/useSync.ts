import { useState, useCallback, useRef, useEffect } from 'react';
import { appApi } from '../services/appApi';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime);

export const useSync = (campaign: string, setStatus: (s: string) => void) => {
    const [localDate, setLocalDate] = useState<string | null>(null);
    const [cloudData, setCloudData] = useState<any>(null);
    const [lastCheck, setLastCheck] = useState<dayjs.Dayjs | null>(null);
    const [pullStatus, setPullStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
    const [isRefreshing, setIsRefreshing] = useState(false);

    const campaignRef = useRef(campaign);
    useEffect(() => {
        campaignRef.current = campaign;
    }, [campaign]);

    const checkStatus = useCallback(async (overrideCampaign?: string) => {
        const campaignId = overrideCampaign || campaignRef.current; // Utilises override when provided, otherwise current campaign state
        if (!campaignId) return;
        
        setIsRefreshing(true);
        setCloudData(null);
        try {
            const local = await appApi.getLocalSaveStatus();
            const cloud = await appApi.getCloudSaveStatus();
            setLocalDate(local === "Never" ? null : local);
            setCloudData(cloud);
            setLastCheck(dayjs());
        } catch (e) {
            console.error("Sync check failed", e);
            setLastCheck(dayjs());
        } finally {
            setIsRefreshing(false); // On disable after completion
        }
    }, []);

    const performPull = async () => {
        setPullStatus("loading");
        setStatus("Downloading...");
        try {
            const result = await appApi.downloadLatestSave();
            if (result === "Success") {
                setPullStatus("success");
                setStatus("Watching");
                await checkStatus(); // Refresh dates
                setTimeout(() => setPullStatus("idle"), 3000);
            } else {
                setPullStatus("error");
                setTimeout(() => setPullStatus("idle"), 5000);
            }
        } catch (err) {
            setPullStatus("error");
            setTimeout(() => setPullStatus("idle"), 5000);
        }
    };

    return { localDate, cloudData, lastCheck, pullStatus, isRefreshing, checkStatus, performPull };
};