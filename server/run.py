from flask import Flask, request, render_template
import time

app = Flask(__name__)
name = 'taro'


@app.route("/", methods=['GET'])
def get_index():
    return render_template('index.html', name=name)

@app.route("/slow", methods=['GET'])
def post_index():
    
    time.sleep(3)
    return render_template('index.html', name=name)

if __name__ == "__main__":
    app.run(
        threaded=False,
        debug=True,
        host="0.0.0.0",
        port=5000)