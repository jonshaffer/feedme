# Social Media Feed Aggregator - Template Repository

This repository serves as a template for creating a personalized, aggregated feed of your social media content. By forking this repository and configuring it, you can automatically fetch posts from selected platforms and have them published as a static `feed.json` file on your GitHub Pages. This enables easy consumption of your social media presence from a single, consistent endpoint.

## Getting Started

To get started with your own aggregated social media feed, follow these steps:

### 1. Forking the Repository

Fork this repository to your own GitHub account. This will create a copy of the project under your control, allowing you to customize it.

### 2. Configuring `config.yaml`

After forking, you will need to create a `config.yaml` file in the `config/` directory. This file tells the aggregator which social media feeds to enable and provides any non-sensitive configuration, as well as global output settings.

A `config.yaml.example` file is provided for reference. Copy its content and modify it:

```bash
cp config/config.yaml.example config/config.yaml
```

Edit `config/config.yaml` to enable or disable feeds, provide specific configurations (e.g., RSS feed URLs), and adjust output settings.

Example `config.yaml`:

```yaml
# Global output settings
output_limit: 0 # Set to a positive integer to limit the total number of items in the output. 0 or negative means no limit.
page_size: 0    # Set to a positive integer to enable pagination. Each page will contain this many items. 0 or 1 means no pagination (single file output).
generate_individual_item_files: false # Set to true to generate a separate JSON file for each feed item.
generate_platform_feeds: false      # Set to true to generate separate JSON files for each social media platform.

feeds:
  linkedin:
    enabled: true
  threads:
    enabled: false # Set to true to enable Threads
  x:
    enabled: true
  instagram:
    enabled: false # Set to true to enable Instagram
  reddit:
    enabled: true
  rss:
    enabled: true
    urls:
      - "https://www.example.com/my-blog-feed.xml"
      - "https://www.another-site.com/news.rss"
```

### 3. Setting Up GitHub Secrets

For platforms requiring authentication (LinkedIn, Threads, X, Instagram, Reddit, Strava, Goodreads, Credly), you must store your API keys and tokens as GitHub Secrets in your forked repository. This ensures sensitive information is not exposed in your public repository.

1.  Go to your forked repository on GitHub.
2.  Navigate to `Settings` > `Secrets and variables` > `Actions`.
3.  Click `New repository secret`.
4.  Create secrets with the following names (replace `YOUR_` with the actual platform name, e.g., `LINKEDIN_API_KEY`):
    *   `LINKEDIN_API_KEY`
    *   `LINKEDIN_API_SECRET`
    *   `THREADS_API_KEY`
    *   `THREADS_API_SECRET`
    *   `X_API_KEY`
    *   `X_API_SECRET`
    *   `INSTAGRAM_API_KEY`
    *   `INSTAGRAM_API_SECRET`
    *   `REDDIT_CLIENT_ID`
    *   `REDDIT_CLIENT_SECRET`
    *   `REDDIT_USERNAME`
    *   `REDDIT_PASSWORD`
    *   `STRAVA_CLIENT_ID`
    *   `STRAVA_CLIENT_SECRET`
    *   `GOODREADS_API_KEY`
    *   `GOODREADS_API_SECRET`
    *   `CREDLY_API_KEY`
    *   `CREDLY_API_SECRET`

    (Note: Specific API key/secret names might vary based on the actual API requirements. Refer to the respective platform's developer documentation for exact requirements.)

The Go application will read these secrets as environment variables during its execution within the GitHub Actions workflow.

### 4. Adjusting the Update Schedule

The feed is automatically updated and published to GitHub Pages via a GitHub Actions workflow (`.github/workflows/build-and-publish.yml`). By default, this workflow runs daily at midnight UTC.

To change the update frequency:

1.  Go to your forked repository on GitHub.
2.  Navigate to `Actions`.
3.  Click on the `Build and Publish Feed` workflow.
4.  Click on the three dots (`...`) next to the workflow run and select `View workflow file`.
5.  Click the pencil icon to edit the file.
6.  Locate the `on: schedule:` section and modify the `cron` expression.

Example (to run every 6 hours):

```yaml
on:
  schedule:
    - cron: '0 */6 * * *' # Runs every 6 hours
```

Refer to the GitHub Actions documentation for cron syntax details.

## Supported Feeds

Currently, the aggregator supports fetching data from:

*   LinkedIn
*   Threads
*   X (formerly Twitter)
*   Instagram
*   Reddit
*   Strava
*   Goodreads
*   Credly
*   RSS Feeds
