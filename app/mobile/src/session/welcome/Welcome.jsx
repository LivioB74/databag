import { useEffect } from 'react';
import { Image, View, TouchableOpacity, Text } from 'react-native';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';
import { styles } from './Welcome.styled';
import { Colors } from 'constants/Colors';
import Ionicons from 'react-native-vector-icons/AntDesign';

import session from 'images/session.png';

export function Welcome() {

  return (
    <SafeAreaView style={styles.container}>
      <Text style={styles.header}>Databag</Text>
      <Text style={styles.label}>Communication for the Decentralized Web</Text>
      <Image style={styles.image} source={session} />
      <View style={styles.steps}>
        <Text style={styles.stepstext}>Setup your profile</Text>
        <Ionicons name={'right'} size={18} color={Colors.text} />
        <Text style={styles.stepstext}>Connect with people</Text>
        <Ionicons name={'right'} size={18} color={Colors.text} />
        <Text style={styles.stepstext}>Start a conversation</Text>
      </View>
    </SafeAreaView>
  );
}

