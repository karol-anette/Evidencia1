#include <iostream>
#include <vector>
#include <algorithm>
using namespace std;

int medianaDeMedianes(vector<int>& arr, int k) {
    if (arr.size() <= 5) {
        sort(arr.begin(), arr.end());
        return arr[k];
    }
    vector<vector<int>> subgrupos;
    for (size_t i = 0; i < arr.size(); i += 5) {
        size_t fin = min(i + 5, arr.size());
        vector<int> grupo(arr.begin() + i, arr.begin() + fin);
        subgrupos.push_back(grupo);
    }
    vector<int> medianas;
    for (auto& grupo : subgrupos) {
        sort(grupo.begin(), grupo.end());
        medianas.push_back(grupo[grupo.size() / 2]);
    }
    int mediana = medianaDeMedianes(medianas, medianas.size() / 2);
    vector<int> menores, iguales, mayores;
    for (int x : arr) {
        if (x < mediana) menores.push_back(x);
        else if (x == mediana) iguales.push_back(x);
        else mayores.push_back(x);
    }
    if (k < (int)menores.size())
        return medianaDeMedianes(menores, k);
    else if (k < (int)menores.size() + (int)iguales.size())
        return mediana;
    else
        return medianaDeMedianes(mayores, k - (int)menores.size() - (int)iguales.size());
}

int main() {
    vector<int> arr = {12, 3, 5, 7, 4, 19, 26};
    int k = 3;
    cout << "El k-esimo elemento mas pequeño es: " << medianaDeMedianes(arr, k) << endl;
    return 0;
}
