import matplotlib.pyplot as plt

find = [[], []]
insert = [[], []]
delete = [[], []]

for data in ((find, 'find'), (insert, 'insert'), (delete, 'delete')):
    for line in open(f'result/{data[1]}.csv').readlines():
        if line == '':
            continue

        count, dur = line.split(',')
        count = int(count)
        dur = int(dur)
        data[0][0].append(count)
        data[0][1].append(dur)

print(find)
print(insert)
print(delete)

plt.grid()
plt.ylabel('Время в нс')
plt.xlabel('Количество элементов')

plt.title('Время выполнения операций в 2-3-tree')
plt.plot(find[0], find[1], label='find')
plt.plot(insert[0], insert[1], label='insert')
plt.plot(delete[0], delete[1], label='delete')
plt.legend()
plt.savefig('result/plot.png')