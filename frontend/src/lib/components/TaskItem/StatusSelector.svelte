<script>
    import { onMount } from 'svelte';

    export let value;
    export let onChange;
    export let immediate = true; 
    export let fullWidth = false; 
    

    let selectorContainer;
    let tempValue = value; 

    function handleChange(newValue) {
        if (immediate) {
            onChange(newValue);
        } else {
            tempValue = newValue;
        }
    }

    // Reset temp value when value prop changes
    $: tempValue = value;

    // Expose method to get current temp value
    export function getCurrentValue() {
        return tempValue;
    }

    function handleMouseMove(event) {
        if (!selectorContainer) return;
        
        const rect = selectorContainer.getBoundingClientRect();
        const x = ((event.clientX - rect.left) / rect.width) * 100;
        const y = ((event.clientY - rect.top) / rect.height) * 100;
        
        selectorContainer.style.setProperty('--mouse-x', `${x}%`);
        selectorContainer.style.setProperty('--mouse-y', `${y}%`);
    }

    onMount(() => {
        if (selectorContainer) {
            selectorContainer.addEventListener('mousemove', handleMouseMove);
        }

        return () => {
            if (selectorContainer) {
                selectorContainer.removeEventListener('mousemove', handleMouseMove);
            }
        };
    });
</script>

<div class="status-selector" class:full-width={fullWidth}>
    <div class="selector-container" bind:this={selectorContainer}>
        <select 
            value={tempValue}
            on:change={(e) => handleChange(e.target.value)}
        >
            <option value="pending">PENDING::0x01</option>
            <option value="in_progress">IN_PROGRESS::0x02</option>
            <option value="completed">COMPLETED::0x03</option>
        </select>
        <div class="selector-decoration">
            <span class="status-indicator"></span>
            <svg class="selector-arrow" viewBox="0 0 24 24" width="16" height="16">
                <path 
                    d="M7 10l5 5 5-5" 
                    fill="none" 
                    stroke="currentColor" 
                    stroke-width="2" 
                    stroke-linecap="round"
                />
            </svg>
        </div>
        <div class="selector-glow"></div>
    </div>
</div>


<style>
    .full-width {
        width: 100%;
    }

    .full-width .selector-container {
        width: 100%;
    }

    .full-width select {
        width: 100%;
    }

    .status-selector {
        position: relative;
    }

    .selector-container {
        position: relative;
        display: flex;
        align-items: center;
    }

    select {
        appearance: none;
        background: #2d3436;
        border: 1px solid #0984e3;
        color: #0984e3;
        padding: 0.5rem 2.5rem 0.5rem 1rem;
        border-radius: 3px;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.8rem;
        cursor: pointer;
        transition: all 0.3s ease;
        position: relative;
        z-index: 2;
    }

    /* Status-specific colors */
    select option[value="pending"] {
        color: #ffd32a;
    }

    select option[value="in_progress"] {
        color: #0984e3;
    }

    select option[value="completed"] {
        color: #00b894;
    }

    /* Decoration elements */
    .selector-decoration {
        position: absolute;
        right: 0.5rem;
        top: 50%;
        transform: translateY(-50%);
        display: flex;
        align-items: center;
        gap: 0.2rem;
        pointer-events: none;
        z-index: 3;
    }

    .status-indicator {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: currentColor;
        animation: pulse 1s infinite;
    }

    .selector-arrow {
        transition: transform 0.3s ease;
    }

    /* Glow effect */
    .selector-glow {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: radial-gradient(
            circle at var(--x, 50%) var(--y, 50%),
            rgba(9, 132, 227, 0.1),
            transparent 50%
        );
        opacity: 0;
        transition: opacity 0.3s ease;
        pointer-events: none;
        z-index: 1;
    }

    /* Hover states */
    .selector-container:hover select {
        border-color: #00b894;
        box-shadow: 0 0 10px rgba(0, 184, 148, 0.2);
    }

    .selector-container:hover .selector-glow {
        opacity: 1;
    }

    .selector-container:hover .selector-arrow {
        transform: translateY(2px);
    }

    /* Status-specific styling */
    .selector-container:has(select[value="pending"]) {
        --status-color: #ffd32a;
    }

    .selector-container:has(select[value="in_progress"]) {
        --status-color: #0984e3;
    }

    .selector-container:has(select[value="completed"]) {
        --status-color: #00b894;
    }

    /* Animations */
    @keyframes pulse {
        0%, 100% {
            transform: scale(1);
            opacity: 0.8;
        }
        50% {
            transform: scale(1);
            opacity: 0.2;
        }
    }

    /* Focus styles */
    select:focus {
        outline: none;
        border-color: #00b894;
        box-shadow: 0 0 10px rgba(0, 184, 148, 0.2);
    }

    /* Add scanline effect */
    .selector-container::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 1px;
        background: currentColor;
        opacity: 0.2;
        animation: scan 2s linear infinite;
        pointer-events: none;
    }

    @keyframes scan {
        from { transform: translateY(0); }
        to { transform: translateY(100%); }
    }

    /* Mouse tracking glow effect */
    .selector-container:hover {
        --x: var(--mouse-x, 50%);
        --y: var(--mouse-y, 50%);
    }

    @media screen and (max-width: 400px) {
        .selector-container select {
            font-size: 0.5rem;
        }
    }
</style>