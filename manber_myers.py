class SubstrRank:
    def __init__(self, left_rank=0, right_rank=0, index=0):
        self.left_rank = left_rank
        self.right_rank = right_rank
        self.index = index

def make_ranks(substr_rank, n):
    r = 1
    rank = [-1] * n
    rank[substr_rank[0].index] = r
    for i in range(1, n):
        if (substr_rank[i].left_rank != substr_rank[i-1].left_rank or
			substr_rank[i].right_rank != substr_rank[i-1].right_rank):
            r += 1
        rank[substr_rank[i].index] = r
    return rank

def suffix_array(T):
    n = len(T)
    substr_rank = []

    for i in range(n):
        substr_rank.append(SubstrRank(ord(T[i]), ord(T[i + 1]) if i < n-1 else 0, i))

    substr_rank.sort(key=lambda sr: (sr.left_rank, sr.right_rank))

    l = 2
    while l < n:
        rank = make_ranks(substr_rank, n)

        # CREAR NUEVA LISTA temporal en lugar de modificar los índices originales
        new_substr_rank = []
        for i in range(n):
            new_substr_rank.append(SubstrRank(
                rank[i], 
                rank[i + l] if i + l < n else 0, 
                i  # Mantener el índice actual
            ))
        
        substr_rank = new_substr_rank
        substr_rank.sort(key=lambda sr: (sr.left_rank, sr.right_rank))
        l *= 2

    SA = [substr_rank[i].index for i in range(n)]

    return SA

SA = suffix_array("mississippi")
print(SA)
expected = [10, 7, 4, 1, 0, 9, 8, 6, 3, 5, 2]
print("Esperado: ", expected)
print("¿Correcto?:", SA == expected)

# IMPLEMENTACIÓN DE BWT Y FM-INDEX PARA BÚSQUEDA

def build_bwt(text, suffix_array):
    n = len(text)
    bwt = []
    for i in range(n):
        #caracter  BWT
        pos = suffix_array[i] - 1
        bwt_char = text[pos] if pos >= 0 else text[-1]
        bwt.append(bwt_char)
    return ''.join(bwt)

def build_fm_index(bwt):
    first_col = sorted(bwt)
    #tabla de ocurrencias y conteos 
    chars = sorted(set(bwt))
    occ_table = {char: [0] * (len(bwt) + 1) for char in chars}
    c_table = {char: 0 for char in chars}
    
    #tabla de ocurrencias 
    for i, char in enumerate(bwt):
        for c in chars:
            occ_table[c][i + 1] = occ_table[c][i] + (1 if bwt[i] == c else 0)
    #cuántos caracteres son menores que cada char
    total = 0
    for char in chars:
        c_table[char] = total
        total += occ_table[char][len(bwt)]
    return first_col, occ_table, c_table