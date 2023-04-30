import os
from bottle import route, run, template
import random

@route('/alias/<name>')
def alias(name):
    name_list = list(name)
    name_list[random.randint(0, (len(name_list) - 1))] = "T"
    new_name = str(name_list).replace('[', "").replace(']', "").replace("'", "").replace(",", "").replace(" ", "")
    return {"alias": new_name}

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 8080))
    run(host='localhost', port=port, debug=True)