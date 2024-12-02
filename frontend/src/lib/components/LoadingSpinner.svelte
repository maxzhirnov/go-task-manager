<script>
    import { mdiLoading } from '@mdi/js';
</script>

<div class="loading-container">
    <svg class="spinner" viewBox="0 0 24 24" width="24" height="24">
        <path fill="currentColor" d={mdiLoading} />
    </svg>
</div>

<style>
     .loading-container {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 200px;
    }

    .spinner {
        animation: spin 1s linear infinite;
        width: 40px;
        height: 40px;
        color: #666;
    }

    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }
</style>