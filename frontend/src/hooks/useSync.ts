import { useState, useCallback, useRef, useEffect } from 'react';
import { GetLocalSaveStatus, GetCloudSaveStatus, DownloadLatestSave } from "../../wailsjs/go/application/App";
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
        const campaignId = overrideCampaign || campaignRef.current; // Utilise l'override si présent (démarrage), sinon l'état
        if (!campaignId) return;
        
        setIsRefreshing(true);
        setCloudData(null);
        try {
            const local = await GetLocalSaveStatus();
            const cloud = await GetCloudSaveStatus();
            setLocalDate(local === "Never" ? null : local);
            setCloudData(cloud);
            setLastCheck(dayjs());
        } catch (e) {
            console.error("Sync check failed", e);
            setLastCheck(dayjs());
        } finally {
            setIsRefreshing(false); // On désactive à la fin
        }
    }, []);

    const performPull = async () => {
        setPullStatus("loading");
        setStatus("Downloading...");
        try {
            const result = await DownloadLatestSave();
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