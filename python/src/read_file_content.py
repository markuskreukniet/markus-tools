import pdfplumber

def read_pdf_with_pdfplumber(file_path):
    content = []

    with pdfplumber.open(file_path) as pdf:
        for page in pdf.pages:
            content.append(page.extract_text())

    return "\n".join(content)

# Usage
file_path = "C:\\Users\\testUser\\Desktop\\test\\test.pdf"
text = read_pdf_with_pdfplumber(file_path)
print(text)
