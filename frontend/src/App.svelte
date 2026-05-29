<script>
    import MapView from './MapView.svelte'
    import Lifelist from './Lifelist.svelte'
    let activeTab = 'map'
</script>

<div class="app">   
    <aside class = "sidebar">
        <button
        class = "nav-btn"
        class:active = {activeTab==='map'}
        on:click={()=> activeTab = 'map'}
        >
        Карта
        </button>

        <button
        class = "nav-btn"
        class:active = {activeTab==='lifelist'}
        on:click={()=> activeTab = 'lifelist'}
        >
        Лайфлист
        </button>
    </aside>

    <main class = "content">
         <div class="content-view" style="display: {activeTab === 'map'? 'block' : 'none'}">    
            <MapView {activeTab} />
         </div>
         <div class="content-view" style="display: {activeTab === 'lifelist'? 'block' : 'none'}">
            <Lifelist />
         </div>
    </main>
</div>

<style>
  :global(body, html) {
    margin: 0;
    padding: 0;
    height: 100%;
    width: 100%;
    overflow: hidden;
    font-family: 'Inter', -apple-system, sans-serif;
  }

  .app {
    display: flex;
    height: 100vh;
    width: 100vw;
    background-color: #f8fafc;
  }

  .sidebar {
    width: 200px;
    background-color: #bbf3dd;
    display: flex;
    flex-direction: column;
    padding: 15px;
    gap: 8px;
    box-shadow: 2px 0 5px rgba(0,0,0,0.05);
    z-index: 10;
  }

  .nav-btn {
    all: unset;  
    padding: 12px 16px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 500;
    color: #2c7a4d;
    transition: all 0.2s ease;
    text-align: left;
  }

  .nav-btn:hover {
    background-color: rgba(255, 255, 255, 0.5);
  }

  .nav-btn.active {
    background-color: #456b8a;
    color: white;
  }

  .content {
    flex-grow: 1;
    position: relative;  
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .content-view {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: none; 
  }

  .content-view.active {
    display: block;
  }
</style>