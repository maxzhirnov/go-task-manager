<script>
    import { mdiCheckCircle, mdiClockOutline, mdiDragVertical, mdiChartBar } from '@mdi/js';
</script>

<div class="landing-container">
    <header class="hero">
        <h1>Task Manager</h1>
        <p class="tagline">Streamline your productivity with our intuitive task management solution</p>
        <div class="cta-buttons">
            <a href="/register" class="btn btn-primary">Get Started Free</a>
            <a href="/login" class="btn btn-secondary">Sign In</a>
        </div>
    </header>

    <section class="features">
        <h2>Why Choose Our Task Manager?</h2>
        <div class="feature-grid">
            <div class="feature-card">
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="currentColor" d={mdiCheckCircle} />
                </svg>
                <h3>Simple & Intuitive</h3>
                <p>Quick task creation and easy organization with our user-friendly interface</p>
            </div>

            <div class="feature-card">
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="currentColor" d={mdiDragVertical} />
                </svg>
                <h3>Drag & Drop</h3>
                <p>Effortlessly prioritize tasks with drag-and-drop functionality</p>
            </div>

            <div class="feature-card">
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="currentColor" d={mdiClockOutline} />
                </svg>
                <h3>Real-time Updates</h3>
                <p>Stay on top of your tasks with instant status updates</p>
            </div>

            <div class="feature-card">
                <svg viewBox="0 0 24 24" width="24" height="24">
                    <path fill="currentColor" d={mdiChartBar} />
                </svg>
                <h3>Progress Tracking</h3>
                <p>Monitor your productivity with detailed statistics</p>
            </div>
        </div>
    </section>

    <section class="cta">
        <h2>Ready to Get Started?</h2>
        <p>Join thousands of users who have improved their productivity</p>
        <a href="/register" class="btn btn-primary">Create Free Account</a>
    </section>
</div>

<style>
    .landing-container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 20px;
        font-family: Helvetica, Arial, sans-serif;
    }

    .hero {
        text-align: center;
        padding: 60px 20px;
        background: linear-gradient(to right, #f8f9fa, #e9ecef);
        border-radius: 8px;
        margin-bottom: 40px;
    }

    .hero h1 {
        font-size: 3rem;
        margin-bottom: 1rem;
        color: #2c3e50;
    }

    .tagline {
        font-size: 1.5rem;
        color: #6c757d;
        margin-bottom: 2rem;
    }

    .cta-buttons {
        display: flex;
        gap: 1rem;
        justify-content: center;
    }

    .btn {
        padding: 12px 24px;
        border-radius: 4px;
        text-decoration: none;
        font-weight: bold;
        transition: all 0.2s;
    }

    .btn-primary {
        background-color: #4CAF50;
        color: white;
    }

    .btn-primary:hover {
        background-color: #45a049;
    }

    .btn-secondary {
        background-color: white;
        color: #4CAF50;
        border: 2px solid #4CAF50;
    }

    .btn-secondary:hover {
        background-color: #f8f9fa;
    }

    .features {
        padding: 40px 0;
    }

    .features h2 {
        text-align: center;
        margin-bottom: 40px;
        color: #2c3e50;
    }

    .feature-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 30px;
        padding: 20px;
    }

    .feature-card {
        padding: 30px;
        background: white;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        text-align: center;
        transition: transform 0.2s;
    }

    .feature-card:hover {
        transform: translateY(-5px);
    }

    .feature-card svg {
        color: #4CAF50;
        width: 48px;
        height: 48px;
        margin-bottom: 20px;
    }

    .feature-card h3 {
        margin-bottom: 15px;
        color: #2c3e50;
    }

    .feature-card p {
        color: #6c757d;
        line-height: 1.6;
    }

    .cta {
        text-align: center;
        padding: 60px 20px;
        background: linear-gradient(to right, #f8f9fa, #e9ecef);
        border-radius: 8px;
        margin-top: 40px;
    }

    .cta h2 {
        margin-bottom: 1rem;
        color: #2c3e50;
    }

    .cta p {
        margin-bottom: 2rem;
        color: #6c757d;
    }

    @media (max-width: 768px) {
        .hero h1 {
            font-size: 2rem;
        }

        .tagline {
            font-size: 1.2rem;
        }

        .feature-grid {
            grid-template-columns: 1fr;
        }

        .cta-buttons {
            flex-direction: column;
            align-items: center;
        }

        .btn {
            width: 100%;
            max-width: 300px;
            text-align: center;
        }
    }
</style>