<!-- TechIcon.svelte -->
<script>
    export let type = 'edit'; // edit, delete, complete
    export let size = '24';
    
    const icons = {
        edit: "M4 20h16m-2-12l2-2-2-2-2 2 2 2zm0 0l-8 8-4 1 1-4 8-8",
        delete: "M6 7v13a2 2 0 002 2h8a2 2 0 002-2V7m-12 0h12m-12 0H7m10 0h-1M9 3h6a1 1 0 011 1v3H8V4a1 1 0 011-1zm0 8v5m6-5v5",
        complete: "M20 6L9 17l-5-5m1 8h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z",
        // Add more icons as needed
    };
</script>

<svg 
    viewBox="0 0 24 24" 
    width={size} 
    height={size}
    class="tech-icon {type}"
>
    <path
        fill="none" 
        stroke="currentColor" 
        stroke-width="1.5" 
        stroke-linecap="round" 
        stroke-linejoin="round"
        d={icons[type]}
        class="icon-path"
    />
</svg>

<style>
    .tech-icon {
        display: inline-block;
        transition: all 0.3s ease;
    }

    .icon-path {
        transition: all 0.3s ease;
    }

    /* Icon-specific styles */
    .edit .icon-path {
        --icon-color: #0984e3;
    }

    .delete .icon-path {
        --icon-color: #ff6b6b;
    }

    .complete .icon-path {
        --icon-color: #00b894;
    }
</style>