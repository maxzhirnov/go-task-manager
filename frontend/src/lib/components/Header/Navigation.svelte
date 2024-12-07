<script>
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
</script>

<div>
    <div class="terminal-header">
        <span class="terminal-title">NAVIGATION.sys</span>
    </div>
    <nav class="nav-links">
        <a 
            class="nav-link" 
            class:active={$page.url.pathname === '/'} 
            href="/" 
            on:click|preventDefault={() => goto('/tasks')}
        >
            <span class="prompt">>_</span>
            <span class="link-text">TASKS</span>
        </a>
        <a 
            class="nav-link" 
            class:active={$page.url.pathname === '/statistics'} 
            href="/statistics" 
            on:click|preventDefault={() => goto('/statistics')}
        >
            <span class="prompt">>_</span>
            <span class="link-text">STATS</span>
        </a>
    </nav>
</div>

<style>
    .terminal-box {
        background: #2d3436;
        border: 1px solid #0984e3;
        border-radius: 4px;
        overflow: hidden;
    }

    .terminal-header {
        background: #1c1c1c;
        padding: 0.3rem 0.5rem;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        border-bottom: 1px solid rgba(9, 132, 227, 0.2);
    }

    .terminal-dots {
        display: flex;
        gap: 4px;
    }

    .dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: #636e72;
    }

    .terminal-title {
        color: #00b894;
        font-size: 0.6rem;
        letter-spacing: 0.1em;
    }

    .nav-links {
        display: flex;
        gap: 1rem;
        padding: 0.5rem;
    }

    .nav-link {
        text-decoration: none;
        color: #00b894;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.4rem 0.8rem;
        border-radius: 3px;
        transition: all 0.3s ease;
    }

    .nav-link:hover {
        background: rgba(0, 184, 148, 0.1);
    }

    .nav-link.active {
        background: rgba(0, 184, 148, 0.1);
        box-shadow: 0 0 8px rgba(0, 184, 148, 0.2);
    }

    .prompt {
        color: #0984e3;
        font-size: 0.8rem;
    }

    .link-text {
        font-size: 0.8rem;
    }

    @media screen and (max-width: 768px) {
        .nav-links {
            justify-content: center;
        }
    }
</style>