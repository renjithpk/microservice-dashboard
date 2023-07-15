#!/usr/bin/env python3
from flask import Flask, send_from_directory, request

app = Flask(__name__)

@app.route('/<path:filename>')
def serve_file(filename):
    if filename.endswith('.yaml') or filename.endswith('.css'):
        response = send_from_directory('.', filename)
        response.headers['Cache-Control'] = 'no-cache, no-store, must-revalidate'
        response.headers['Pragma'] = 'no-cache'
        response.headers['Expires'] = '0'
        return response
    else:
        return send_from_directory('.', filename)

if __name__ == '__main__':
    app.run(port=8000)
