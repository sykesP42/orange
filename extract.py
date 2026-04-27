import zipfile, os, shutil

html_dir = "/tmp/orange-html"
zip_path = "/tmp/orange-frontend.zip"

os.makedirs(html_dir, exist_ok=True)
for f in os.listdir(html_dir):
    p = os.path.join(html_dir, f)
    if os.path.isdir(p):
        shutil.rmtree(p)
    elif f != 'assets' and f != 'pwa-icons':
        os.remove(p)

z = zipfile.ZipFile(zip_path, 'r')
total = len(z.namelist())
for name in z.namelist():
    fixed = name.replace('\\', '/')
    if fixed.startswith('/'):
        fixed = fixed[1:]
    if fixed.endswith('/'):
        os.makedirs(os.path.join(html_dir, fixed), exist_ok=True)
        continue
    parts = [p for p in fixed.split('/') if p and p != '.']
    target = os.path.join(html_dir, *parts) if parts else None
    if not target:
        continue
    os.makedirs(os.path.dirname(target), exist_ok=True)
    try:
        with z.open(name) as src, open(target, 'wb') as dst:
            dst.write(src.read())
    except Exception as e:
        print(f'Error: {name}: {e}')
print(f'Extracted {total} files to {html_dir}')
z.close()
