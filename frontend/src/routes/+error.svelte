<script>
    import { page } from '$app/stores';

    // Custom messages for different status codes
    $: message = {
        404: "The page you're looking for doesn't exist.",
        500: "Something went wrong on our server.",
        403: "You don't have permission to access this page.",
    }[$page.status] || 'An unexpected error occurred.';
</script>

<div class="error-container">
    <div class="error-content">
        <div class="error-code">{$page.status}</div>
        <h1 class="error-title">Oops! {$page.error?.message || 'Page not found'}</h1>
        <p class="error-message">{message}</p>
        <div class="action-buttons">
            <a href="/" class="button primary">Go Home</a>
            <button class="button secondary" on:click={() => history.back()}>
                Go Back
            </button>
        </div>
    </div>
</div>

<style>
    .error-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: calc(100vh - 64px);
        background: linear-gradient(135deg, #f5f7fa 0%, #e4e8eb 100%);
        padding: 20px;
    }

    .error-content {
        background: white;
        padding: 40px;
        border-radius: 20px;
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        text-align: center;
        max-width: 500px;
        width: 90%;
        animation: slideUp 0.5s ease-out;
    }

    .error-code {
        font-size: 96px;
        font-weight: bold;
        color: #4CAF50;
        line-height: 1;
        margin-bottom: 20px;
        text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
    }

    .error-title {
        font-size: 24px;
        color: #333;
        margin: 0 0 15px 0;
        font-weight: 600;
    }

    .error-message {
        color: #666;
        font-size: 16px;
        line-height: 1.6;
        margin: 0 0 30px 0;
    }

    .action-buttons {
        display: flex;
        gap: 15px;
        justify-content: center;
    }

    .button {
        padding: 12px 24px;
        border-radius: 8px;
        font-size: 16px;
        font-weight: 500;
        text-decoration: none;
        transition: all 0.3s ease;
        cursor: pointer;
        border: none;
    }

    .button.primary {
        background-color: #4CAF50;
        color: white;
        box-shadow: 0 4px 6px rgba(76, 175, 80, 0.2);
    }

    .button.primary:hover {
        background-color: #45a049;
        transform: translateY(-2px);
        box-shadow: 0 6px 8px rgba(76, 175, 80, 0.3);
    }

    .button.secondary {
        background-color: #f5f5f5;
        color: #333;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .button.secondary:hover {
        background-color: #e8e8e8;
        transform: translateY(-2px);
        box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
    }

    @keyframes slideUp {
        from {
            opacity: 0;
            transform: translateY(20px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    /* Responsive adjustments */
    @media (max-width: 480px) {
        .error-code {
            font-size: 72px;
        }

        .error-title {
            font-size: 20px;
        }

        .action-buttons {
            flex-direction: column;
        }

        .button {
            width: 100%;
        }
    }
</style>