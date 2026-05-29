<script>
  import { onMount, tick } from 'svelte';
  import { GetLifelist, SearchBirds, AddBird, ImportLifelist} from '../wailsjs/go/main/App';

  let searchText = '';
  let searchResults = [];
  let myBirds = [];
  
  let isLoading = true; 
  let errorMessage = '';

  onMount(loadLifelist);

  async function loadLifelist() {
    try {
      const data = await GetLifelist();
      myBirds = data || [];
    } catch (e) {
      console.error("Ошибка загрузки:", e);
      errorMessage = "Не удалось загрузить данные";
    } finally {
      isLoading = false;
    }
  }

  async function handleSearch() {
    if (searchText.trim().length < 2) {
      searchResults = [];
      return;
    }
    try {
      searchResults = await SearchBirds(searchText) || [];
    } catch (e) {
      console.error("Ошибка поиска:", e);
    }
  }

  async function handleImport() {
    try {
      await ImportLifelist();
      
      await loadLifelist(); 
      alert("Импорт успешно завершен!");
    } catch (e) {
      console.error("Ошибка импорта:", e);
      alert("Произошла ошибка при импорте файла.");
    }
  }

  async function addBird(bird) {
    try {
      await AddBird(bird);
      myBirds = [bird, ...myBirds];
      searchText = '';
      searchResults = [];
    } catch (e) {
      alert("Ошибка при сохранении");
    }
  }
</script>

<div class="lifelist-wrapper">
  {#if isLoading}
    <div class="status-msg"> Загрузка данных из системы...</div>
  {:else if errorMessage}
    <div class="status-msg error"> {errorMessage}</div>
  {:else}
    <div class="content">
      
      <header>
        <div class="title-area">
          <h2> Мой Лайфлист</h2>
          <span class="counter">Всего видов: {myBirds.length}</span>
        </div>
        
        <button class="import-btn" on:click={handleImport}>
           Импорт CSV
        </button>
      </header>

      <div class="search-box">
        <input 
          type="text" 
          bind:value={searchText} 
          on:input={handleSearch}
          placeholder="Начните вводить название птицы..."
        />
        
        {#if searchResults.length > 0}
          <div class="dropdown">
            {#each searchResults as bird}
              <div class="item" on:click={() => addBird(bird)}>
                <!-- Если приходят английские названия, выводим их -->
                <strong>{bird.comName}</strong> 
                <small>{bird.sciName}</small>
              </div>
            {/each}
          </div>
        {/if}
      </div>

      <div class="birds-list">
        <h3>Наблюдения ({myBirds.length})</h3>
        {#each myBirds as bird}
          <div class="bird-card">
            <span>{bird.comName}</span>
            <span class="latin">{bird.sciName}</span>
          </div>
        {:else}
          <p>Список пуст. Добавьте птицу через поиск.</p>
        {/each}
      </div>
    </div>
  {/if}
</div>

<style>
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    border-bottom: 1px solid #e2e8f0;
    padding-bottom: 12px;
  }

  .title-area {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .title-area h2 {
    margin: 0;
  }

  .counter {
    background: #e2f5ee;
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 0.9em;
    color: #2c7a4d;
    font-weight: bold;
  }

  .import-btn {
    background-color: #ffffff;
    color: #456b8a;
    border: 1px solid #456b8a;
    padding: 8px 16px;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  .lifelist-wrapper { padding: 20px; 
    color: #333;
    height: 100%;           
    overflow-y: auto;      
    box-sizing: border-box;   }
  .status-msg { padding: 20px; text-align: center; background: #f0f0f0; border-radius: 8px; }
  .error { color: red; background: #ffdada; }
  
  .search-box { position: relative; margin-bottom: 20px; }
  input { width: 100%; padding: 12px; border: 1px solid #ccc; border-radius: 6px; }
  
  .dropdown { 
    position: absolute; top: 50px; left: 0; right: 0; 
    background: white; border: 1px solid #ddd; z-index: 1000;
    max-height: 200px; overflow-y: auto; box-shadow: 0 4px 10px rgba(0,0,0,0.1);
  }
  .item { padding: 10px; 
          cursor: pointer; 
          border-bottom: 1px solid #eee; 
          display: flex; 
          flex-direction: column; 
        }
  .item:hover { background: #e8f4ff; }
  
  .bird-card { 
    padding: 10px; 
    background: white; 
    border-bottom: 1px solid #eee;
    display: flex; 
    justify-content: space-between;
  }
  .latin { color: #888;
   font-style: italic; 
   font-size: 0.9em; }
</style>