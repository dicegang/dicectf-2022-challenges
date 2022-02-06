from Crypto.Util.number import long_to_bytes
# flag = b"dice{w0rld_of_w1res_03294803}"
# key = b'V]\xd2\xd7\xc6\xcc\xa8KZVp\xd8\xd5p\xcc\xa9\xcbW\xca\xf0\xa8*\xab%\xae$(*\xc3'
flag = b"dice"

def get_key(flag):
    q = [int(a) for a in "".join(format(i, '08b') for i in flag)]

    stage1 = [q[0]]

    for i in range(len(q) - 1):
        stage1.append(q[i + 1] ^ q[i])

    key = long_to_bytes(int("".join([str(i) for i in stage1]), 2))


    stage2 = [int(a) for a in "".join(format(i, '08b') for i in key)]



    for i in range(len(stage2)):
        stage2[i] = stage2[i] ^ stage1[i]
    return key
