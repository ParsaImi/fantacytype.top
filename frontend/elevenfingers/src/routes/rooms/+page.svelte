<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  
  let ws: WebSocket;
  let rooms: Record<string, Record<string, any>> = {};
  let loading = true;
  let error = null;
  let userNickname: string = '';
  let isAuthenticated: boolean = false;
  let username = '';
  let token = '';
  
  onMount(() => {
    token = localStorage.getItem('auth_token');
    if (token) {
      try {
        // Basic validation and extract username from token
        const parts = token.split('.');
        if (parts.length === 3) {
          // Decode the payload (middle part of the JWT)
          const payload = JSON.parse(atob(parts[1]));
          if (payload && payload.sub) {
            username = payload.sub;
            userNickname = username
            isAuthenticated = true;
          }
        }
      } catch (error) {
        console.error('Error decoding token:', error);
      }
    } else {
        userNickname = localStorage.getItem('user_nickname') || '';
        
        if (!userNickname) {
          // Redirect to nickname page if no nickname is set
          goto('/nickname?redirect=/game');
          return;
        }
    } 
    
    // Connect to WebSocket server
    ws = new WebSocket('wss://wss.fantacytype.top/ws');
    
    ws.onopen = () => {
      console.log('Connected to the server');
      // Request the list of rooms
      requestRooms();
    };
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      switch(data.type) {
        case 'roomsStatus':
          rooms = data.rooms || {};
          loading = false;
          break;
        default:
          console.log('Received unknown message type:', data.type);
      }
    };
    
    ws.onerror = (err) => {
      console.error('WebSocket error:', err);
      error = 'Failed to connect to the server. Please try again later.';
      loading = false;
    };
    
    ws.onclose = () => {
      console.log('Disconnected from the server');
    };
  });
  
  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
  
  function requestRooms() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'roomsStatus'
      }));
    }
  }
  
  function refreshRooms() {
    loading = true;
    requestRooms();
  }
  
  function joinRoom(roomName: string) {
    // Store the room name to be used in the game page
    localStorage.setItem('selected_room', roomName);
    // Navigate to the game page
    goto('/game');
  }
  
  function createNewRoom() {
    // Open a modal or prompt for room name
    const roomName = prompt('Enter a new room name:');
    if (roomName && roomName.trim()) {
      localStorage.setItem('selected_room', roomName.trim());
      goto('/game');
    }
  }
  
  function changeNickname() {
    goto('/nickname?redirect=/rooms');
  }
</script>

<svelte:head>
  <title>Available Rooms - Tenfinger Typing Game</title>
  <meta name="description" content="Join a multiplayer typing game room" />
</svelte:head>

<div class="container">
  {#if !isAuthenticated}
    <div class="user-bar">
      <div class="user-info">
        <span>Playing as: <strong>{userNickname}</strong></span>
      </div>
      <button class="change-nickname" on:click={changeNickname}>
        Change Nickname
      </button>
    </div>
  {/if}
  
  <div class="rooms-container">
    <h2>Available Rooms</h2>
    
    {#if loading}
      <div class="loading">
        <p>Loading available rooms...</p>
        <div class="loading-spinner"></div>
      </div>
    {:else if error}
      <div class="error">
        <p>{error}</p>
        <button class="btn-retry" on:click={refreshRooms}>Try Again</button>
      </div>
    {:else}
      <div class="room-controls">
        <button class="btn-secondary" on:click={refreshRooms}>
          Refresh Rooms
        </button>
        <button class="btn-primary" on:click={createNewRoom}>
          Create New Room
        </button>
      </div>
      
      {#if Object.keys(rooms).length === 0}
        <div class="no-rooms">
          <p>No active rooms found. Create a new room to start playing!</p>
        </div>
      {:else}
        <div class="rooms-list">
          {#each Object.entries(rooms) as [roomName, players]}
            <div class="room-card">
              <div class="room-info">
                <h3>{roomName}</h3>
                <p>Players: {Object.keys(players).length}</p>
              </div>
              <button class="btn-join" on:click={() => joinRoom(roomName)}>
                Join Room
              </button>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  /* Dark Mode Styling - Modern UI for Rooms Page */
  :global(body) {
    font-family: 'IRANSans', 'Vazir', 'Tahoma', 'Segoe UI', -apple-system, BlinkMacSystemFont, sans-serif;
    background-color: #121212;
    color: #e0e0e0;
    margin: 0;
    padding: 0;
    min-height: 100vh;
    transition: background-color 0.3s ease;
  }
  
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
  }
  
  .user-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    background-color: #1e1e1e;
    border-radius: 8px;
    margin-bottom: 2rem;
    border: 1px solid #333;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  }
  
  .user-info {
    font-size: 1.1rem;
    color: #e0e0e0;
  }
  
  .user-info strong {
    color: #bb86fc;
  }
  
  .change-nickname {
    padding: 0.5rem 1rem;
    background-color: transparent;
    border: 1px solid #bb86fc;
    color: #bb86fc;
    border-radius: 4px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s, color 0.2s;
  }
  
  .change-nickname:hover {
    background-color: #bb86fc;
    color: #121212;
    box-shadow: 0 2px 4px rgba(187, 134, 252, 0.3);
  }
  
  .rooms-container {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  h2 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    text-align: center;
    color: #bb86fc;
    text-shadow: 0 2px 8px rgba(187, 134, 252, 0.3);
  }
  
  .loading, .error, .no-rooms {
    width: 100%;
    max-width: 600px;
    background-color: #1e1e1e;
    padding: 2rem;
    border-radius: 12px;
    text-align: center;
    margin-bottom: 2rem;
    border: 1px solid #333;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  }
  
  .loading-spinner {
    margin: 20px auto;
    width: 40px;
    height: 40px;
    border: 4px solid rgba(187, 134, 252, 0.3);
    border-radius: 50%;
    border-top: 4px solid #bb86fc;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  .error {
    background-color: #1e1e1e;
    border: 1px solid #cf6679;
  }
  
  .loading p, .error p, .no-rooms p {
    margin-bottom: 1rem;
    color: #aaa;
    font-size: 1.1rem;
  }
  
  .room-controls {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
    width: 100%;
    max-width: 600px;
    justify-content: space-between;
  }
  
  .btn-primary, .btn-secondary, .btn-retry, .btn-join {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  }
  
  .btn-secondary {
    background-color: #2d2d2d;
    color: #e0e0e0;
    border: 1px solid #424242;
  }
  
  .btn-secondary:hover {
    background-color: #3d3d3d;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  .btn-primary {
    background-color: #bb86fc;
    color: #121212;
    border: none;
  }
  
  .btn-primary:hover {
    background-color: #c9a0ff;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.3);
  }
  
  .btn-retry {
    background-color: #2d2d2d;
    color: #e0e0e0;
    border: 1px solid #424242;
    margin-top: 1rem;
  }
  
  .btn-retry:hover {
    background-color: #3d3d3d;
  }
  
  .rooms-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    width: 100%;
    max-width: 600px;
  }
  
  .room-card {
    background-color: #1e1e1e;
    border-radius: 12px;
    padding: 1.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    border: 1px solid #333;
  }
  
  .room-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
    border-color: #bb86fc;
  }
  
  .room-info h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1.25rem;
    color: #bb86fc;
  }
  
  .room-info p {
    margin: 0;
    color: #aaa;
  }
  
  .btn-join {
    padding: 0.5rem 1rem;
    background-color: #bb86fc;
    color: #121212;
    border: none;
    font-weight: 600;
  }
  
  .btn-join:hover {
    background-color: #c9a0ff;
    transform: translateY(-2px);
    box-shadow: 0 2px 4px rgba(187, 134, 252, 0.3);
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    h2 {
      font-size: 2rem;
    }
    
    .room-card {
      flex-direction: column;
      text-align: center;
      gap: 1rem;
    }
    
    .room-controls {
      flex-direction: column;
      gap: 0.75rem;
    }
    
    .btn-primary, .btn-secondary {
      width: 100%;
    }
  }
  
  @media (max-width: 480px) {
    .user-bar {
      flex-direction: column;
      gap: 1rem;
    }
    
    .container {
      padding: 1rem;
    }
  }
</style>
