from playwright.sync_api import sync_playwright
import os

def run_cuj(page):
    print("Navigating to http://localhost:8080")
    try:
        page.goto("http://localhost:8080", timeout=10000)
    except Exception as e:
        print(f"Failed to navigate: {e}")
        return

    page.wait_for_timeout(2000)
    print("Page title:", page.title())

    # Take screenshot
    page.screenshot(path="verification_result.png")
    print("Screenshot saved to verification_result.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(viewport={"width": 1280, "height": 720})
        page = context.new_page()
        try:
            run_cuj(page)
        finally:
            context.close()
            browser.close()
