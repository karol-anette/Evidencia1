import time
import sys
import tracemalloc #Herramienta que les permita perfilar el uso de memoria

def getBuckets(T):
    count = {}
    buckets = {}
    for c in T:
        count[c] = count.get(c, 0) + 1
    start = 0
    for c in sorted(count.keys()):
        buckets[c] = (start, start + count[c])
        start += count[c]
    return buckets


def place_LMS(SA, LMS, count, buckets, T, t): #Función para colocar LMS
    end = None
    for i in range(len(T) - 1, 0, -1):
        if t[i] == "S" and t[i - 1] == "L":
            revoffset = count[T[i]] = count.get(T[i], 0) + 1
            SA[buckets[T[i]][1] - revoffset] = i
            if end is not None:
                LMS[i] = end
            end = i
    LMS[len(T) - 1] = len(T) - 1

def induce_L(SA, count, buckets, T, t):
    for i in range(len(SA)):
        if SA[i] >= 0 and t[SA[i] - 1] == "L":
            symbol = T[SA[i] - 1]
            offset = count.get(symbol, 0)
            SA[buckets[symbol][0] + offset] = SA[i] - 1
            count[symbol] = offset + 1

def induce_S(SA, count, buckets, T, t):
    for i in range(len(SA) - 1, 0, -1):
        if SA[i] > 0 and t[SA[i] - 1] == "S":
            symbol = T[SA[i] - 1]
            revoffset = count[symbol] = count.get(symbol, 0) + 1
            SA[buckets[symbol][1] - revoffset] = SA[i] - 1

def sais(T):
    t = ["_"] * len(T)

    t[- 1] = "S"
    for i in range(len(T) - 1, 0, -1):
        if T[i-1] < T[i]:
            t[i - 1] = "S"
        elif T[i - 1] == T[i] and t[i] == "S":
            t[i - 1] = "S"
        else:
            t[i - 1] = "L"
    
    buckets = getBuckets(T)

    count = {}
    SA = [-1] * len(T)
    LMS = {}
  
    place_LMS(SA, LMS, count, buckets, T, t)

    count.clear()
    induce_L(SA, count, buckets, T, t)

    count.clear()
    induce_S(SA, count, buckets, T, t)
 
    namesp = [-1] * len(T)
    name = 0
    prev = None
    LMS_indices = [i for i in range(len(T)) if t[i] == "S" and t[i - 1] == "L"]
    for i in LMS_indices:
        end_i = LMS.get(i, len(T))
        if prev is not None:
            end_prev = LMS.get(prev, len(T))
            if T[prev:end_prev] != T[i:end_i]:
                name += 1
        namesp[i] = name
        prev = i

    names = []
    SApIdx = []
    for i in range(len(T)):
        if namesp[i] != -1:
            names.append(namesp[i])
            SApIdx.append(i)

    if name < len(names) - 1:
        names = sais(names)

    SA = [-1] * len(T)
    count.clear()
    for i, pos in enumerate(SApIdx):
        idx = names[i]
        revoffset = count[T[pos]] = count.get(T[pos], 0) + 1
        SA[buckets[T[pos]][1] - revoffset] = pos

    count.clear()
    induce_L(SA, count, buckets, T, t)
    count.clear()
    induce_S(SA, count, buckets, T, t)

    return SA
#BWT
def build_bwt(text, SA):
    n = len(SA)
    bwt = []
    for i in SA:
        if i == 0:
            bwt.append('$')
        else:
            bwt.append(text[i - 1])
    return ''.join(bwt)

def build_fm_index(bwt):
    alphabet = sorted(set(bwt))
    C = {}
    total = 0
    for c in alphabet:
        C[c] = total
        total += bwt.count(c)
    Occ = {c: [0] for c in alphabet}
    for ch in bwt:
        for c in alphabet:
            Occ[c].append(Occ[c][-1] + (1 if ch == c else 0))
    return C, Occ

def fm_search(pattern, C, Occ, bwt):
    # IMplementamos backward search.
    l, r = 0, len(bwt)
    for c in reversed(pattern):
        if c not in C:
            return []
        l = C[c] + Occ[c][l]
        r = C[c] + Occ[c][r]
        if l >= r:
            return []
    return list(range(l, r))

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Elige un archivo")
        sys.exit(1)

    filename = sys.argv[1]

    with open(filename, "r", encoding="utf-8") as f:
        text = f.read().strip() + "$" 

    T = list(text.lower())
    
    tracemalloc.start()
    start_time = time.time()

    SA = sais(T)
    # Construimos BWT y FM-index.
    bwt = build_bwt(T, SA)
    C, Occ = build_fm_index(bwt)
    
    # Buscar una cadena exacta
    pattern = input("Ingresa la cadena a buscar: ").lower()
    text_lower = text.lower()
    bwt_lower = bwt.lower()
    C_lower, Occ_lower = build_fm_index(bwt_lower)

    matches = []
    i = 0
    while True:
        i = text_lower.find(pattern, i)
        if i == -1:
            break
        matches.append(i)
        i += 1 

    if matches:
        print(f"La cadena '{pattern}' aparece {len(matches)} veces.")
    else:
        print(f"No se encontraron resultados para '{pattern}'.")


    end_time = time.time()
    current, peak = tracemalloc.get_traced_memory()
    tracemalloc.stop()

    print(f"Tiempo de ejecución: {end_time - start_time:.6f} segundos.")
    print(f"Uso máximo de memoria: {peak / 1024 / 1024:.2f} MB") #Le damos display a memory profiling.
    print(f"Tamañoo: {len(T)} caracteres")
    #print("Primeros 100 índices:", SA[:100])