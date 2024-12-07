<!-- TechButton.svelte -->
<script>
    import TechIcon from './TechIcon.svelte';
    
    export let type = 'edit';
    export let title = '';
    export let disabled = false;
</script>

<button
    class="tech-button {type}"
    {title}
    aria-label={title}
    {disabled}
    on:click
>
    <TechIcon {type} />
    <div class="button-glow"></div>
</button>

<style>
    .tech-button {
        background: transparent;
        border: 1px solid currentColor;
        color: #0984e3;
        padding: 0.4rem;
        border-radius: 3px;
        cursor: pointer;
        transition: all 0.3s ease;
        position: relative;
        overflow: hidden;
    }

    .button-glow {
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(9, 132, 227, 0.2),
            transparent
        );
        transition: 0.5s;
    }

    .tech-button:hover .button-glow {
        left: 100%;
    }

    /* Button type styles */
    .tech-button.edit {
        --button-color: #0984e3;
        --button-hover-color: #00b894;
        --button-glow-color: rgba(0, 184, 148, 0.3);
    }

    .tech-button.delete {
        --button-color: #ff6b6b;
        --button-hover-color: #ff4757;
        --button-glow-color: rgba(255, 107, 107, 0.3);
    }

    .tech-button.complete {
        --button-color: #00b894;
        --button-hover-color: #00d8a8;
        --button-glow-color: rgba(0, 184, 148, 0.3);
    }

    .tech-button:hover {
        color: var(--button-hover-color);
        box-shadow: 0 0 8px var(--button-glow-color);
    }

    .tech-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    /* Reduced motion */
    @media (prefers-reduced-motion: reduce) {
        .button-glow {
            display: none;
        }
    }

    @media screen and (max-width: 400px) {
        .tech-button {
            padding: 0.2rem;
        }
    }
</style>