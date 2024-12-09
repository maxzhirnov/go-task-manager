<script>
    import { browser } from "$app/environment";

    export let description;
    export let taskId;
    
    const MAX_CHARS = browser && window.innerWidth < 700 ? 30 : 60; // Maximum characters to show when collapsed
    const isLongDescription = description.length > MAX_CHARS || description.includes('\n');
    
    // Default to collapsed state for long descriptions
    const initialState = browser ? 
        JSON.parse(localStorage.getItem(`task_desc_${taskId}`)) ?? true : 
        true;
    
    let isCollapsed = initialState && isLongDescription;

    function toggleDescription() {
        if (!isLongDescription) return;
        isCollapsed = !isCollapsed;
        if (browser) {
            localStorage.setItem(`task_desc_${taskId}`, JSON.stringify(isCollapsed));
        }
    }

    function getTruncatedText(text, maxLength) {
        if (text.length <= maxLength && !text.includes('\n')) return text;
        
        // If there's a newline before maxLength, truncate there
        const newlineIndex = text.indexOf('\n');
        if (newlineIndex > -1 && newlineIndex < maxLength) {
            return text.substring(0, newlineIndex);
        }
        
        // Otherwise truncate at last space before maxLength
        const truncated = text.substring(0, maxLength);
        const lastSpace = truncated.lastIndexOf(' ');
        return truncated.substring(0, lastSpace > 0 ? lastSpace : maxLength);
    }
</script>

<div class="task-text-content">
    <p class:collapsed={isCollapsed}>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <span class="desc-label" class:clickable={isLongDescription} on:click={toggleDescription}>
            [DESC{isLongDescription ? (isCollapsed ? ' +' : ' -') : ''}]
        </span>
        {#if isCollapsed}
            {getTruncatedText(description, MAX_CHARS)}...
        {:else}
            {description}
        {/if}
    </p>
</div>

<style>
    .task-text-content p {
        white-space: pre-line;
        font-family: "JetBrains Mono", monospace;
        font-size: 0.9rem;
        line-height: 1.6;
        color: #fff;
        background: #2d3436;
        padding: 1rem;
        border-radius: 3px;
        /* border-left: 2px solid #0984e3; */
        position: relative;
        margin: 0.5rem 0;
        transition: all 0.3s ease;
    }

    .collapsed {
        padding: 0.5rem 1rem !important;
        min-height: 1.5rem;
        /* border-left: 2px solid #0984e3; */
        background: #2d3436;
    }

    .desc-label {
        position: absolute;
        top: -0.7rem;
        left: 0.0rem;
        font-size: .6rem;
        color: #00b894;
        background: #2d3436;
        padding: 0.1rem 0.3rem;
        letter-spacing: 0.1em;
        user-select: none;
        z-index: 1;
    }

    .clickable {
        cursor: pointer;
    }

    .clickable:hover {
        color: #00d1a7;
    }

    @media screen and (max-width: 400px) {
        .task-text-content p {
            font-size: 0.6rem;
        }
    }
</style>