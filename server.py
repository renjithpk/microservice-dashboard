#!/usr/bin/env python3
from flask import Flask, send_from_directory, request, jsonify

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
@app.route('/details/booking-api/<path:path>')
def test_endpoint_1(path):
    # Create a simple JSON response
    data = {
        'message': f'This is a test endpoint with path: {path}',
        'status': 'success',
        'api': {
            'version': "v1.2.3"
        }
    }
    # Use jsonify to convert the dictionary to JSON format and return it as the response
    return jsonify(data)
@app.route('/details/booking-api/<path:path>')
def test_endpoint_2(path):
    # Create a simple JSON response
    data = {
        'message': f'This is a test endpoint with path: {path}',
        'status': 'success',
        'api': {
            'version': "v0.2.3"
        }
    }
    # Use jsonify to convert the dictionary to JSON format and return it as the response
    return jsonify(data)

if __name__ == '__main__':
    app.run(port=8000)
