<script>
    import { onMount, tick } from 'svelte'; // 1. Добавили tick
    import L from 'leaflet';
    import 'leaflet/dist/leaflet.css';
    import { GetLocations, GetLifelist } from '../wailsjs/go/main/App';

    import markerIcon from 'leaflet/dist/images/marker-icon.png';
    import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png';
    import markerShadow from 'leaflet/dist/images/marker-shadow.png';

    const myDefaultIcon = L.icon({
        iconUrl: markerIcon,
        iconRetinaUrl: markerIcon2x,
        shadowUrl: markerShadow,
        iconSize: [25, 41],         // Стандартный размер иконки Leaflet
        iconAnchor: [12, 41],       // Точка иконки, которая смотрит на координату
        popupAnchor: [1, -34],      // Откуда будет "вырастать" попап
        shadowSize: [41, 41]        // Размер тени
    });

    L.Marker.prototype.options.icon = myDefaultIcon;
    
    export let activeTab; 

    let mapContainer;
    let map;
    let markerGroup;  
    
    let locations = [];
    let myLifelistNames = new Set();  
    let isLoading = true;
    
    let showLocations = true;
    let showOnlyNewBirds = false;

    async function getIpLocation() {
        try {
            const res = await fetch("http://ip-api.com/json/");
            const data = await res.json();
            if (data.status === "success") return { lat: data.lat, lng: data.lon };
            throw new Error("API fail");
        } catch {
            return { lat: 53.9, lng: 27.5 }; 
        }
    }

    
    function normalizeName(name) {
        return name ? name.toLowerCase().trim() : '';
    }

    onMount(async () => {
        let pos = await getIpLocation();
        map = L.map(mapContainer, { zoomControl: false }).setView([pos.lat, pos.lng], 10);
        
        L.control.zoom({ position: 'bottomright' }).addTo(map);

        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; OpenStreetMap contributors'
        }).addTo(map);

        markerGroup = L.layerGroup().addTo(map);

        try {
            const [locData, lifelistData] = await Promise.all([
                GetLocations(),
                GetLifelist()
            ]);
            
            locations = locData || [];
            if (lifelistData) {
                myLifelistNames = new Set(lifelistData.map(bird => normalizeName(bird.comName)));
            }
        } catch (e) {
            console.error("Ошибка загрузки данных для карты:", e);
        } finally {
            isLoading = false; 
        }
    });

    $: if (activeTab === 'map' && map) {
        handleTabSwitch();
    }   

    $: if (map && locations) {
        renderMarkers(showLocations, showOnlyNewBirds);
    }

   async function handleTabSwitch() {
    await tick(); 

    try {
        const lifelistData = await GetLifelist();
        
        myLifelistNames = new Set((lifelistData || []).map(bird => normalizeName(bird.comName)));

        renderMarkers(showLocations, showOnlyNewBirds);

        setTimeout(() => {
            if (map) map.invalidateSize();
        }, 100);

    } catch (e) {
        if (map) map.invalidateSize();
    }
}

    function renderMarkers(showLocs, showNew) {
        if (!markerGroup) return;
        
        markerGroup.clearLayers();

        if (!showLocs) return;

        locations.forEach(location => {
            let validObservations = location.observations || [];

            if (showNew) {
                validObservations = validObservations.filter(obs => !myLifelistNames.has(normalizeName(obs.comName)));
            }

            if (showNew && validObservations.length === 0) return;

            let marker = L.marker([location.lat, location.lng]);

            let popupHtml = `
                <div class="custom-popup">
                    <div class="popup-title">${location.locName}</div>
            `;
            
            if (validObservations.length > 0) {
                popupHtml += `<ul class="popup-list">`;
                validObservations.forEach(obs => {
                    popupHtml += `
                        <li>
                            <span class="bird-name">${obs.comName}</span>
                            <span class="bird-date">${obs.obsDt}</span>
                        </li>`;
                });
                popupHtml += `</ul>`;
            } else {
                popupHtml += `<p class="no-data">Деталей по птицам нет</p>`;
            }
            
            popupHtml += `</div>`;
            
            marker.bindPopup(popupHtml);
            markerGroup.addLayer(marker);
        });
    }
</script>

<div class="map-wrapper">

    {#if isLoading}
        <div class="loader-overlay">
            <div class="spinner"></div>
            <p class="loader-text">Ищем птиц поблизости...</p>
        </div>
    {/if}
    <div bind:this={mapContainer} class="map-container"></div>

    <div class="controls-panel">
        <h3 class="panel-title">Настройки карты</h3>
        
        <label class="toggle-row">
            <input type="checkbox" bind:checked={showLocations} />
            <span class="toggle-text">Показывать локации</span>
        </label>

        <label class="toggle-row">
            <input type="checkbox" bind:checked={showOnlyNewBirds} disabled={!showLocations}/>
            <span class="toggle-text">Только птицы не из лайфлиста</span>
        </label>
    </div>
</div>

<style>
    .map-wrapper {
        flex: 1;
        position: relative;
        height: 100vh;
        width: 100%;
        display: flex;
        flex-direction: column;
    }

    .map-container {
        height: 100%;
        width: 100%;
        z-index: 1;  
    }

    .controls-panel {
        position: absolute;
        top: 20px;
        right: 20px;
        z-index: 1000; 
        background: rgba(255, 255, 255, 0.95);
        backdrop-filter: blur(5px);
        padding: 16px 20px;
        border-radius: 12px;
        box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
        border: 1px solid rgba(0,0,0,0.05);
        min-width: 220px;
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .panel-title {
        margin: 0 0 4px 0;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
        font-size: 14px;
        color: #666;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .toggle-row {
        display: flex;
        align-items: center;
        gap: 10px;
        cursor: pointer;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
        font-size: 15px;
        color: #333;
        transition: opacity 0.2s;
    }

    .toggle-row input {
        width: 18px;
        height: 18px;
        cursor: pointer;
        accent-color: #2c7a4d;  
    }

    .toggle-row input:disabled + .toggle-text {
        opacity: 0.5;
        text-decoration: line-through;
    }

    :global(.custom-popup) {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        max-height: 250px;
        overflow-y: auto;
    }

    :global(.popup-title) {
        font-weight: bold;
        font-size: 15px;
        color: #2c7a4d;
        border-bottom: 2px solid #eee;
        padding-bottom: 6px;
        margin-bottom: 8px;
    }

    :global(.popup-list) {
        list-style: none;
        padding: 0;
        margin: 0;
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    :global(.popup-list li) {
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-bottom: 1px dashed #eee;
        padding-bottom: 4px;
    }

    :global(.bird-name) {
        font-weight: 500;
        color: #333;
    }

    :global(.bird-date) {
        font-size: 11px;
        color: #888;
        margin-left: 10px;
    }

    :global(.no-data) {
        color: #999;
        font-style: italic;
        margin: 0;
    }
    .loader-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(255, 255, 255, 0.7);
        backdrop-filter: blur(4px); 
        z-index: 2000; 
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }

    .spinner {
        width: 48px;
        height: 48px;
        border: 4px solid #e2e8f0;
        border-top: 4px solid #2c7a4d;  
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin-bottom: 16px;
    }

    .loader-text {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
        color: #2c7a4d;
        font-weight: 500;
        font-size: 16px;
        margin: 0;
    }

    @keyframes spin {
        0% { transform: rotate(0deg); }
        100% { transform: rotate(360deg); }
    }
</style>