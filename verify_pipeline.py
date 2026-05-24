from playwright.sync_api import sync_playwright

def run_cuj(page):
    # Navigate to the registry
    page.goto("http://localhost:8080")
    page.wait_for_timeout(2000)

    # Click Access Console to enter the new pipeline UI
    # We use a broader selector or text search
    try:
        page.get_by_text("Access Console").first.click()
        page.wait_for_timeout(2000)
        print("Entered Project Console")
    except Exception as e:
        print(f"Could not click Access Console: {e}")
        # If no projects, the button won't exist.
        return

    # Now we should be in the Project Console with the sub-menu
    # Take a screenshot of the new UI
    page.screenshot(path="pipeline_ui.png")

    # Navigate through tabs to show they work for the video
    tabs = ["Spatial Alignment", "Wall Analytics & Inspection", "AR Field Sync", "Data Ingestion"]
    for tab in tabs:
        page.get_by_text(tab).click()
        page.wait_for_timeout(1000)
        page.screenshot(path=f"tab_{tab.replace(' ', '_').lower()}.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(viewport={"width": 1920, "height": 1080}, record_video_dir="videos")
        page = context.new_page()
        try:
            run_cuj(page)
        finally:
            context.close()
            browser.close()
