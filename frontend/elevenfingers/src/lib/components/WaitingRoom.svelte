<script lang="ts">
  export let players: Record<string, boolean> = {};
  export let onReady: () => void;
  
  $: readyCount = Object.values(players).filter(ready => ready).length;
  $: totalPlayers = Object.keys(players).length;
</script>

<div class="waiting-room">
  <h2>Waiting Room</h2>
  
  <div class="player-list">
    <h3>Players ({readyCount}/{totalPlayers} ready)</h3>
    {#if totalPlayers === 0}
      <p>No players in room yet. Waiting for others to join...</p>
    {:else}
      <ul>
        {#each Object.entries(players) as [playerName, isReady]}
          <li class={isReady ? 'ready' : 'not-ready'}>
            {playerName} 
            {#if isReady}
              <span class="ready-icon">✓</span>
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </div>
  
  <div class="ready-section">
    <p>Once you're ready to play, click the button below:</p>
    <button on:click={onReady}>I'm Ready</button>
    <p class="hint">Game will start when all players are ready</p>
  </div>
</div>

<style>
  /* Dark Mode Styling - Modern UI for Waiting Room */
  .waiting-room {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
    color: #e0e0e0;
  }
  
  h2 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    text-align: center;
    color: #bb86fc;
    text-shadow: 0 2px 8px rgba(187, 134, 252, 0.3);
  }
  
  .player-list {
    width: 100%;
    max-width: 500px;
    margin-bottom: 2rem;
    padding: 1.5rem;
    border: 1px solid #333;
    border-radius: 12px;
    background-color: #1e1e1e;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }
  
  .player-list:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
    border-color: #bb86fc;
  }
  
  h3 {
    margin-top: 0;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #333;
    color: #bb86fc;
  }
  
  ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
  }
  
  li {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid #333;
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: background-color 0.2s ease;
  }
  
  li:last-child {
    border-bottom: none;
  }
  
  li:hover {
    background-color: #2d2d2d;
  }
  
  .ready {
    color: #03dac6; /* Teal accent color from Material Design dark theme */
    font-weight: 600;
  }
  
  .not-ready {
    color: #aaa;
  }
  
  .ready-icon {
    display: inline-block;
    margin-left: 8px;
    color: #03dac6;
    animation: pulse 2s infinite;
  }
  
  @keyframes pulse {
    0% { opacity: 0.7; }
    50% { opacity: 1; }
    100% { opacity: 0.7; }
  }
  
  .ready-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    background-color: #1e1e1e;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #333;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    width: 100%;
    max-width: 500px;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }
  
  .ready-section:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
    border-color: #bb86fc;
  }
  
  .ready-section p {
    color: #e0e0e0;
    margin-bottom: 1rem;
  }

  button {
    padding: 12px 28px;
    font-size: 1rem;
    font-weight: 600;
    color: #121212;
    background-color: #bb86fc;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    margin: 1rem 0;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  }
  
  button:hover {
    background-color: #c9a0ff;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  .hint {
    font-style: italic;
    color: #aaa;
    font-size: 0.9rem;
    margin-top: 0.5rem;
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    h2 {
      font-size: 2rem;
    }
    
    .player-list, .ready-section {
      padding: 1.25rem;
    }
  }
  
  @media (max-width: 480px) {
    .waiting-room {
      padding: 1rem;
    }
    
    button {
      width: 100%;
    }
  }
</style>
