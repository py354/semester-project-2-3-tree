from os.path import join
import random
import time
import os

random.seed(time.time())

for v in range(1, 11):
    for count in [100, 1000, 10000, 100000, 200000, 400000, 600000, 800000, 1000000]:
        path = join('dataset/', f'v{v}')
        os.makedirs(path, exist_ok=True)
        with open(join(path, f'{count}.csv'), 'w') as file:
            for _ in range(count):
                file.write(str(random.randint(0, 2**31-1)) + ',')
