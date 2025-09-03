<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import WaitingRoom from '$lib/components/WaitingRoom.svelte';
  import TypingGame from '$lib/components/TypingGame.svelte';
  
  let ws: WebSocket;
  let currentView = 'waiting'; // 'waiting', 'game'
  let roomData = { players: {} };
  let gameData: any = null;
  let typingGameComponent: TypingGame;
  let selectedRoom: string = '';
  let isPersianRoom = false; // Flag for Persian room
  let userNickname: string = '';
  let username = '';
  let token = '';
  let isAuthenticated = false;
  
  onMount(() => {
    // Check if user has a nickname
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
    }else{
        userNickname = localStorage.getItem('user_nickname') || '';
        
        if (!userNickname) {
          // Redirect to nickname page if no nickname is set
          goto('/nickname?redirect=/game');
          return;
        }
    }
    
    
    // Get the selected room from localStorage
    selectedRoom = localStorage.getItem('selected_room') || 'room1';
    
    // Check if it's the Persian room
    isPersianRoom = selectedRoom === 'room3';
    
    // Connect to WebSocket server
    ws = new WebSocket('wss://ws.fantacytype.top/ws'); // Replace with your actual WebSocket URL
    
    ws.onopen = () => {
      console.log('Connected to the server');
        if (token) {

          ws.send(JSON.stringify({
              type: 'usercred',
              content: {
                token : token
              }
          }))

        } else {
          ws.send(JSON.stringify({
              type: 'usercred',
              content: {
                username : userNickname
              }
          }))
        }
      // Join the selected room when the page loads
      joinGame(selectedRoom);
    };
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      switch(data.type) {
        case 'roomStatus':
          roomData = data;
          currentView = 'waiting';
          break;
        case 'startGame':
          gameData = data;
          currentView = 'game';
          break;
        case 'userProgress':
          // Handle progress updates from server
          if (typingGameComponent && data.userid && typeof data.percentage === 'number') {
            typingGameComponent.updateProgress(data.userid, data.percentage);
          }
          break;
        case 'playerRank':
          // Handle player rank updates
            typingGameComponent.updatePlayerRanks(data);
            console.log("HERE IS DDDDDDDDDDDDDAAAAAAAAAAAAAAAAAAAAATTTTTTTTTTTAAAAAAAAAAAAAAA :::::::::::::::::::")
            console.log(data)
          break;
        case 'endGame':
          // Handle game end signal
          if (typingGameComponent) {
            typingGameComponent.endGame(data);
          }
          break;
        default:
          console.log('Received unknown message type:', data.type);
      }
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
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
  
  function joinGame(roomName: string = selectedRoom) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'join',
        content: { 
          room: roomName,
          nickname: userNickname
        }
      }));
    }
  }
  
  function sendReady() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'ready'
      }));
    }
  }
  
  function sendWordComplete(word: string) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'wordComplete',
        content: { word }
      }));
    }
  }
  
  function handlePlayAgain() {
    // Return to waiting room view
    currentView = 'waiting';
    
    // Rejoin the game room
    joinGame();
  }
  
  function changeRoom() {
    // Navigate back to the rooms list
    goto('/rooms');
  }
  
  function changeNickname() {
    goto('/nickname?redirect=/game');
  }
</script>

<svelte:head>
  <title>Tenfinger Typing Game - {selectedRoom} Room</title>
  <meta name="description" content="Multiplayer typing game" />
</svelte:head>

<div class="container {isPersianRoom ? 'persian-container' : ''}">
  {#if !isAuthenticated}
    <div class="user-bar">
      <div class="user-info">
        <span>Playing as: <strong>{userNickname}</strong></span>
      </div>
      <div class="user-controls">
        <button class="change-nickname" on:click={changeNickname}>
          Change Nickname
        </button>
      </div>
    </div>
  {/if}

  {#if currentView === 'waiting'}
    <div class="room-header">
      <h2>Room: {selectedRoom} {isPersianRoom ? '(Persian)' : ''}</h2>
      <button class="change-room-btn" on:click={changeRoom}>Change Room</button>
    </div>
    <WaitingRoom 
      players={roomData.players} 
      onReady={sendReady}
      isPersian={isPersianRoom}
      currentUserNickname={userNickname}
    />
  {:else if currentView === 'game'}
    <TypingGame 
      gameData={gameData} 
      onWordComplete={sendWordComplete}
      bind:this={typingGameComponent}
      on:playAgain={handlePlayAgain}
      isPersian={isPersianRoom}
      currentUserNickname={userNickname}
    />
  {/if}
</div>

<style>
  /* Dark Mode Styling to match TypingGame component */
  :global(body) {
    font-family: 'IRANSans', 'Vazir', 'Tahoma', 'Segoe UI', -apple-system, BlinkMacSystemFont, sans-serif;
    background-color: #121212;
    color: #e0e0e0;
    text-align: center;
    margin: 0;
    padding: 0;
    min-height: 100vh;
    transition: background-color 0.3s ease;
  }
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .persian-container {
    direction: rtl;
    text-align: right;
  }
  
  /* User bar styling */
  .user-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    background-color: #1e1e1e;
    border-radius: 12px;
    margin-bottom: 1.5rem;
    width: 100%;
    max-width: 800px;
    border: 1px solid #333;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  }
  
  .user-info {
    font-size: 1.1rem;
    color: #e0e0e0;
  }
  
  .user-info strong {
    color: #bb86fc;
    font-weight: 600;
  }
  
  .change-nickname {
    padding: 0.6rem 1.2rem;
    background-color: #2d2d2d;
    border: 1px solid #bb86fc;
    color: #bb86fc;
    border-radius: 8px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
  
  .change-nickname:hover {
    background-color: #bb86fc;
    color: #121212;
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  }
  
  /* Room header styling */
  .room-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid #333;
    width: 100%;
    max-width: 800px;
  }
  
  .room-header h2 {
    margin: 0;
    font-size: 2rem;
    color: #bb86fc;
    font-weight: 600;
    letter-spacing: -0.5px;
  }
  
  .change-room-btn {
    padding: 0.6rem 1.2rem;
    background-color: #2d2d2d;
    border: 1px solid #424242;
    color: #e0e0e0;
    border-radius: 8px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
  
  .change-room-btn:hover {
    background-color: #3d3d3d;
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    .container {
      padding: 1rem;
    }
    
    .room-header {
      flex-direction: column;
      gap: 1rem;
      text-align: center;
    }
    
    .room-header h2 {
      font-size: 1.5rem;
    }
  }
  
  @media (max-width: 480px) {
    .user-bar {
      flex-direction: column;
      gap: 1rem;
    }
    
    .change-nickname, .change-room-btn {
      width: 100%;
    }
  }
</style>
