import csv

provincias = []
cantones = []
distritos = []


def capitalizador(inputString):
    lista = inputString.split(" ")
    outputString = ""
    for i in lista:
        outputString = outputString+i.capitalize()+" "
    outputString = outputString.rstrip()
    # print(outputString)
    return (outputString)


with open("distelec0.csv", newline='') as distelec:
    lector = csv.reader(distelec)
    for row in lector:
        # print(row)
        # print(row[0])
        newProvincia = [row[0][:-5], capitalizador(row[1])]
        if newProvincia not in provincias:
            provincias.append(newProvincia)
        newCanton = [row[0][:-3], capitalizador(row[2])]
        if newCanton not in cantones:
            cantones.append(newCanton)
        newDistrito = [row[0], capitalizador(row[3])]
        # if newDis not in provincias:
        distritos.append(newDistrito)

# print(provincias)
# print(cantones)
# print(distritos)

for i in provincias:
    print(f'INSERT INTO Provincias VALUES ({i[0]},"{i[1]}");')

for i in cantones:
    print(f'INSERT INTO Cantones VALUES ({i[0]},"{i[1]}",{i[0][:-2]});')

for i in distritos:
    print(f'INSERT INTO Distritos VALUES ({i[0]},"{i[1]}",{i[0][:-3]});')
