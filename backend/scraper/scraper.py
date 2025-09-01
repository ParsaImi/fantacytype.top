import requests
import sqlite3
import time
import re
from bs4 import BeautifulSoup
from urllib.parse import urljoin
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent
DB_PATH = BASE_DIR / "data" / "real_typing_sentences_fa.db"

def setup_db():
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    cursor.execute('''
               CREATE TABLE IF NOT EXISTS typing_sentences (
                   id INTEGER PRIMARY KEY AUTOINCREMENT,
                   content TEXT NOT NULL,
                   language TEXT
                   )
               ''')
    return conn, cursor


def extract_links_from_main_page(url):
    try:
        response = requests.get(url)
        homeSoup = BeautifulSoup(response.text)
        a_tags = homeSoup.select("div.main-content li.report h3 a, div.main-content li.news h3 a")
        links = []
        for a_tag in a_tags:
            href = a_tag.get('href')
            if href:
                # Convert relative URLs to absolute URLs
                absolute_url = urljoin(url, href)
                links.append(absolute_url)

        return links

    except requests.RequestException as e:
        print(f"Error fetching main page: {e}")
        return []



def scrape_article_content(url, cursor):
    try:
        r = requests.get(url)
        r.raise_for_status()
        soup = BeautifulSoup(r.text)
        pees = soup.select("#item p")
        pattern = re.compile(r'^[آابپتثجچحخدذرزژسشصضطظعغفقکگلمنوهیئ\u200C\u200D., ،]+$')
        article_added = 0
        for item in pees:
            if 'همشهری' not in item.text:
                if len(item.text) > 120 and len(item.text) < 330:
                    if pattern.match(item.text):
                        new_text = item.text.replace('‌', ' ')
                        new_text = new_text.replace('،', '')
                        cursor.execute(
                                "INSERT INTO typing_sentences (content, language) VALUES (?, ?)",
                                (new_text , "fa_IR")
                        )
                        print(new_text)
                        print("PPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPP")
                        article_added += 1
        print(f"Added {article_added} from {url}")
        return article_added
    except requests.RequestException as e:
        print(f"Error while scraping : {e}")
        return 0











def main():
    print("started")
    conn, cursor = setup_db()
    main_page_urls = ["https://www.hamshahrionline.ir/service/Lifeskills/heathsubpage",
                     "https://www.hamshahrionline.ir/service/Lifeskills/%D8%AF%D9%83%D9%88%D8%B1%D8%A7%D8%B3%D9%8A%D9%88%D9%86",
                     "https://www.hamshahrionline.ir/service/Lifeskills/tips",
                     "https://www.hamshahrionline.ir/service/Society/highereducation",
                     "https://www.hamshahrionline.ir/service/Science/Hi-Tech",
                     "https://www.hamshahrionline.ir/service/Sport/soccer-world",
                     "https://www.hamshahrionline.ir/service/Sport/ballandnet"
                     ]


    try:
        total_sentences = 0
        all_article_links = []
        for main_page_url in main_page_urls:
            print("ONE MAIN PAGEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
            article_links = extract_links_from_main_page(main_page_url)
            all_article_links.extend(article_links)
            time.sleep(1)

        all_article_links = list(set(all_article_links))
        for i, link in enumerate(all_article_links, 1):
            print("DONNNNNNNNNNNE TIMEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
            print(f"\nProcessing article {i}/{len(all_article_links)}")
            sentences_added = scrape_article_content(link, cursor)
            total_sentences += sentences_added
            
            # Commit after each article to prevent data loss
            conn.commit()
            
            # Add delay to be respectful to the server
            time.sleep(1)
        print(f"\nScraping completed! Total sentences added: {total_sentences}")
    except Exception as e:
        print(f"Error in main function : {e}")
    finally:
        conn.close

if __name__ == "__main__":
    main()
