"""Make the business examples used by the Forms demo, independently of test fixtures.

Run with the bundled Python (Pillow). Templates are authored by the adjacent
generate_xlsx.py, generate_office.py and generate_pdfme.py scripts.
"""
import json
from pathlib import Path

from PIL import Image, ImageDraw

OUT = Path(__file__).with_name("samples")


def save(name, value):
    (OUT / name).write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    # A small document mark, drawn at twice its displayed size. It is a logo
    # placeholder in the editable example, rather than a test image gradient.
    mark = Image.new("RGBA", (160, 160), (0, 0, 0, 0))
    draw = ImageDraw.Draw(mark)
    draw.rounded_rectangle((0, 0, 159, 159), radius=24, fill="#263342")
    draw.rounded_rectangle((39, 29, 112, 130), radius=8, fill="#FFFFFF")
    draw.line((55, 59, 96, 59), fill="#263342", width=6)
    draw.line((55, 78, 96, 78), fill="#263342", width=6)
    draw.line((55, 97, 82, 97), fill="#263342", width=6)
    mark.save(OUT / "brand-mark.png", optimize=True)

    save("form-fixed.json", {
        "no": 1024,
        "date": "2026-10-14",
        "customer": {"name": "株式会社みなと商事"},
        "product": "デスクオーガナイザー",
        "amount": 12000,
        "total": 13200,
        "paid": False,
        "memo": "下記内容をご確認ください。\n総額は税込表示です。",
        "photo": "brand-mark.png",
        "url": "https://example.com/orders/1024",
        "jan": "4901234567894",
        "tags": ["オーク材", "幅280 mm", "国内製造"],
    })

    products = [
        ("ノート A5", 10, 480),
        ("ボールペン", 20, 150),
        ("デスクトレー", 3, 2800),
        ("ファイルボックス", 6, 650),
        ("付箋セット", 4, 420),
    ]
    items = [{"name": name, "qty": qty, "price": price, "amount": qty * price} for name, qty, price in products]
    subtotal = sum(item["amount"] for item in items)
    tax = subtotal // 10
    save("form-rows.json", {
        "no": "INV-2026-1014",
        "customer": "株式会社みなと商事",
        "date": "2026-10-14",
        "due": "2026-11-13",
        "subtotal": subtotal,
        "tax": tax,
        "total": subtotal + tax,
        "note": "お振込手数料は貴社にてご負担ください。\n商品お届け後の請求分です。",
        "items": items,
    })

    delivery = []
    for i in range(18):
        name, _, price = products[i % len(products)]
        qty = i % 3 + 1
        delivery.append({"no": i + 1, "name": name, "qty": qty, "amount": qty * price})
    pages, carry = [], 0
    for start in range(0, len(delivery), 10):
        rows = delivery[start:start + 10]
        end = carry + sum(row["amount"] for row in rows)
        pages.append({"carry_in": carry, "page_total": end - carry, "carry_out": end, "rows": rows})
        carry = end
    save("form-pages.json", {
        "slip": "DN-2026-1014",
        "customer": "株式会社みなと商事",
        "date": "2026-10-14",
        "total": carry,
        "pages": pages,
    })


if __name__ == "__main__":
    main()
